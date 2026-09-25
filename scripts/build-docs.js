#!/usr/bin/env node
/**
 * Cross-platform documentation build script for LVBP GameCast
 * 
 * Usage:
 *   node scripts/build-docs.js           # Full build (default)
 *   node scripts/build-docs.js --ci      # CI mode (exit on error, no dev server)
 *   node scripts/build-docs.js --dev     # Dev mode (build + start dev server)
 *   node scripts/build-docs.js --help    # Show help
 */

import { spawnSync, spawn } from 'child_process';
import { existsSync, mkdirSync } from 'fs';
import { join, resolve } from 'path';
import { fileURLToPath } from 'url';

const __dirname = resolve(fileURLToPath(import.meta.url), '..');
const ROOT = resolve(__dirname, '..');
const BACKEND_DIR = join(ROOT, 'backend');
const WEB_DIR = join(ROOT, 'web');
const DOCS_DIR = join(ROOT, 'docs-site');
const DOCS_PUBLIC_API = join(DOCS_DIR, 'public', 'api');
const DOCS_PUBLIC_TYPEDOC = join(DOCS_DIR, 'public', 'typedoc');

function parseArgs() {
  const args = process.argv.slice(2);
  return {
    ci: args.includes('--ci'),
    dev: args.includes('--dev'),
    help: args.includes('--help') || args.includes('-h'),
    skipOpenapi: args.includes('--skip-openapi'),
    skipTypedoc: args.includes('--skip-typedoc'),
    skipVitepress: args.includes('--skip-vitepress'),
  };
}

function printHelp() {
  console.log(`
╔═══════════════════════════════════════════════════════════════════╗
║           LVBP GameCast - Documentation Builder                   ║
╠═══════════════════════════════════════════════════════════════════╣
║  Usage: node scripts/build-docs.js [options]                      ║
╠═══════════════════════════════════════════════════════════════════╣
║  Options:                                                         ║
║    --ci              CI mode: exit on error, no dev server        ║
║    --dev             Build and start VitePress dev server         ║
║    --skip-openapi    Skip OpenAPI spec generation (swag)          ║
║    --skip-typedoc    Skip TypeScript docs generation (typedoc)    ║
║    --skip-vitepress  Skip VitePress build                         ║
║    --help, -h        Show this help                               ║
╠═══════════════════════════════════════════════════════════════════╣
║  Examples:                                                        ║
║    npm run build-docs              # Full production build        ║
║    npm run build-docs:ci           # CI/CD pipeline build         ║
║    npm run dev-docs                # Build + dev server           ║
║    node scripts/build-docs.js --skip-openapi --skip-typedoc       ║
╚═══════════════════════════════════════════════════════════════════╝
`);
}

function log(stage, message, isError = false) {
  const prefix = isError ? '❌' : '✅';
  const timestamp = new Date().toISOString().split('T')[1].slice(0, -1);
  console.log(`[${timestamp}] ${prefix} [${stage}] ${message}`);
}

function runCommand(cmd, args, cwd, stage) {
  log(stage, `Running: ${cmd} ${args.join(' ')}`);
  
  const result = spawnSync(cmd, args, {
    cwd,
    stdio: 'inherit',
    shell: true,
    env: { ...process.env, FORCE_COLOR: '1' },
  });

  if (result.error) {
    log(stage, `Command failed to start: ${result.error.message}`, true);
    return false;
  }

  if (result.status !== 0) {
    log(stage, `Exited with code ${result.status}`, true);
    return false;
  }

  log(stage, 'Completed successfully');
  return true;
}

function ensureDir(dir) {
  if (!existsSync(dir)) {
    mkdirSync(dir, { recursive: true });
    log('SETUP', `Created directory: ${dir}`);
  }
}

async function generateOpenAPI() {
  log('OPENAPI', 'Generating OpenAPI spec from Go backend...');
  
  // Check if swag is available
  const swagCheck = spawnSync('go', ['tool', 'github.com/swaggo/swag/cmd/swag', 'version'], {
    cwd: BACKEND_DIR,
    shell: true,
    stdio: 'pipe',
  });

  if (swagCheck.status !== 0) {
    log('OPENAPI', 'Installing swag...', false);
    const installResult = spawnSync('go', ['install', 'github.com/swaggo/swag/cmd/swag@latest'], {
      cwd: BACKEND_DIR,
      shell: true,
      stdio: 'inherit',
    });
    if (installResult.status !== 0) {
      log('OPENAPI', 'Failed to install swag', true);
      return false;
    }
  }

  ensureDir(DOCS_PUBLIC_API);

  return runCommand(
    'go',
    ['tool', 'github.com/swaggo/swag/cmd/swag', 'init', '-g', 'cmd/api/main.go', '-o', DOCS_PUBLIC_API, '--parseDependency', '--parseInternal'],
    BACKEND_DIR,
    'OPENAPI'
  );
}

async function generateTypedoc() {
  log('TYPEDOC', 'Generating TypeScript documentation...');

  // Check if node_modules exists in web
  if (!existsSync(join(WEB_DIR, 'node_modules'))) {
    log('TYPEDOC', 'Installing frontend dependencies...', false);
    const installResult = spawnSync('pnpm', ['install', '--frozen-lockfile'], {
      cwd: WEB_DIR,
      shell: true,
      stdio: 'inherit',
    });
    if (installResult.status !== 0) {
      log('TYPEDOC', 'Failed to install frontend dependencies', true);
      return false;
    }
  }

  ensureDir(DOCS_PUBLIC_TYPEDOC);

  return runCommand(
    'npx',
    ['typedoc', '--out', DOCS_PUBLIC_TYPEDOC, '--entryPoints', 'src/lib', '--plugin', 'typedoc-plugin-markdown', '--readme', 'none'],
    WEB_DIR,
    'TYPEDOC'
  );
}

async function buildVitePress() {
  log('VITEPRESS', 'Building VitePress documentation site...');

  // Ensure docs-site dependencies
  if (!existsSync(join(DOCS_DIR, 'node_modules'))) {
    log('VITEPRESS', 'Installing docs-site dependencies...', false);
    const installResult = spawnSync('pnpm', ['install', '--frozen-lockfile'], {
      cwd: DOCS_DIR,
      shell: true,
      stdio: 'inherit',
    });
    if (installResult.status !== 0) {
      log('VITEPRESS', 'Failed to install docs-site dependencies', true);
      return false;
    }
  }

  return runCommand('pnpm', ['run', 'build'], DOCS_DIR, 'VITEPRESS');
}

function startDevServer() {
  log('DEV', 'Starting VitePress dev server...');
  
  const child = spawn('pnpm', ['run', 'dev', '--', '--host', '0.0.0.0', '--port', '5173'], {
    cwd: DOCS_DIR,
    shell: true,
    stdio: 'inherit',
    env: { ...process.env, FORCE_COLOR: '1' },
  });

  child.on('error', (err) => {
    log('DEV', `Failed to start dev server: ${err.message}`, true);
    process.exit(1);
  });

  // Handle graceful shutdown
  process.on('SIGINT', () => {
    log('DEV', 'Shutting down dev server...');
    child.kill('SIGINT');
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    child.kill('SIGTERM');
    process.exit(0);
  });
}

async function main() {
  const options = parseArgs();

  if (options.help) {
    printHelp();
    process.exit(0);
  }

  console.log('\n🚀 LVBP GameCast Documentation Builder\n');
  log('INIT', `Root: ${ROOT}`);
  log('INIT', `Mode: ${options.ci ? 'CI' : options.dev ? 'DEV' : 'PRODUCTION'}`);

  const startTime = Date.now();
  let allSuccess = true;

  // Step 1: Generate OpenAPI spec
  if (!options.skipOpenapi) {
    const success = await generateOpenAPI();
    allSuccess = allSuccess && success;
    if (!success && options.ci) process.exit(1);
  } else {
    log('OPENAPI', 'Skipped (--skip-openapi)');
  }

  // Step 2: Generate TypeScript docs
  if (!options.skipTypedoc) {
    const success = await generateTypedoc();
    allSuccess = allSuccess && success;
    if (!success && options.ci) process.exit(1);
  } else {
    log('TYPEDOC', 'Skipped (--skip-typedoc)');
  }

  // Step 3: Build VitePress
  if (!options.skipVitepress) {
    const success = await buildVitePress();
    allSuccess = allSuccess && success;
    if (!success && options.ci) process.exit(1);
  } else {
    log('VITEPRESS', 'Skipped (--skip-vitepress)');
  }

  const duration = ((Date.now() - startTime) / 1000).toFixed(1);

  if (allSuccess) {
    log('DONE', `All steps completed in ${duration}s`, false);
    
    if (options.dev) {
      await startDevServer();
    } else {
      console.log('\n📁 Output: docs-site/.vitepress/dist/');
      console.log('🌐 Deploy: Push to main branch for GitHub Pages auto-deploy\n');
    }
  } else {
    log('DONE', `Build completed with errors in ${duration}s`, true);
    if (options.ci) process.exit(1);
  }
}

main().catch((err) => {
  console.error('💥 Unexpected error:', err);
  process.exit(1);
});