#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const COLORS = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m',
  red: '\x1b[31m',
  magenta: '\x1b[35m',
};

function log(color, message) {
  console.log(`${COLORS[color]}${message}${COLORS.reset}`);
}

function findRepoRoot() {
  let current = path.resolve(process.cwd());
  while (current !== path.dirname(current)) {
    const skillMdPath = path.join(current, 'SKILL.md');
    if (fs.existsSync(skillMdPath)) {
      return current;
    }
    const skillsDir = path.join(current, 'skills');
    if (fs.existsSync(skillsDir) && fs.statSync(skillsDir).isDirectory()) {
      const skillFiles = fs.readdirSync(skillsDir).filter(f => f.endsWith('.md') || fs.statSync(path.join(skillsDir, f)).isDirectory());
      if (skillFiles.length > 0) {
        return current;
      }
    }
    current = path.dirname(current);
  }
  return process.cwd();
}

function loadSkills(repoRoot) {
  const skillsDir = path.join(repoRoot, 'skills');
  if (!fs.existsSync(skillsDir)) {
    log('red', 'Error: skills/ directory not found');
    return [];
  }
  const entries = fs.readdirSync(skillsDir);
  const skillDirs = entries.filter(e => {
    const fullPath = path.join(skillsDir, e);
    return fs.statSync(fullPath).isDirectory();
  }).sort();
  const skills = [];
  for (const dir of skillDirs) {
    const skillMdPath = path.join(skillsDir, dir, 'SKILL.md');
    if (fs.existsSync(skillMdPath)) {
      const stat = fs.statSync(skillMdPath);
      if (stat.size > 0) {
        skills.push({ name: dir, path: skillMdPath, size: stat.size });
      } else {
        log('yellow', `Warning: ${dir}/SKILL.md exists but is empty`);
      }
    } else {
      log('yellow', `Warning: ${dir}/SKILL.md not found`);
    }
  }
  const mdFiles = entries.filter(e => e.endsWith('.md') && !fs.statSync(path.join(skillsDir, e)).isDirectory());
  for (const md of mdFiles) {
    const mdPath = path.join(skillsDir, md);
    const stat = fs.statSync(mdPath);
    if (stat.size > 0) {
      skills.push({ name: md, path: mdPath, size: stat.size });
    }
  }
  return skills;
}

function verifyFiles(skills) {
  let allValid = true;
  for (const skill of skills) {
    if (!fs.existsSync(skill.path)) {
      log('red', `Error: File not found: ${skill.path}`);
      allValid = false;
      continue;
    }
    const stat = fs.statSync(skill.path);
    if (stat.size <= 0) {
      log('red', `Error: File is empty: ${skill.path}`);
      allValid = false;
    }
  }
  return allValid;
}

function createSymlink(repoRoot) {
  const linkPath = path.join(repoRoot, 'skills');
  const targetPath = repoRoot;
  if (fs.existsSync(linkPath)) {
    try {
      const stat = fs.lstatSync(linkPath);
      if (stat.isSymbolicLink()) {
        fs.unlinkSync(linkPath);
      } else if (stat.isDirectory()) {
        return;
      }
    } catch (err) {
      // ignore
    }
  }
  try {
    fs.symlinkSync(targetPath, linkPath, 'junction');
    log('green', `Symlink created: skills/ -> ${repoRoot}`);
  } catch (err) {
    log('yellow', `Could not create symlink: ${err.message}`);
  }
}

function printInstructions(repoRoot, skills) {
  console.log('');
  log('bright', '============================================');
  log('bright', '  Pro Skills - Senior Full-Stack Engineer');
  log('bright', '  Setup Complete!');
  log('bright', '============================================');
  console.log('');
  log('cyan', 'Loaded Skills:');
  for (const skill of skills) {
    log('green', `  - ${skill.name} (${skill.size} bytes)`);
  }
  console.log('');
  log('yellow', 'Next Steps:');
  log('blue', `  1. cd ${repoRoot}`);
  log('blue', '  2. npm run setup');
  log('blue', '  3. Explore skills/ directory');
  log('blue', '  4. Read each SKILL.md for details');
  console.log('');
  log('magenta', `  Repository: ${repoRoot}`);
  log('magenta', '  Author: Eng. Salah Allsayani');
  log('magenta', '  License: MIT');
  console.log('');
}

function main() {
  log('bright', 'Setting up Pro Skills...');
  const repoRoot = findRepoRoot();
  log('green', `Repository root found: ${repoRoot}`);
  const skills = loadSkills(repoRoot);
  if (skills.length === 0) {
    log('red', 'No skills found!');
    process.exit(1);
  }
  log('green', `Loaded ${skills.length} skill(s)`);
  const allValid = verifyFiles(skills);
  if (!allValid) {
    log('red', 'Some files failed verification');
    process.exit(1);
  }
  createSymlink(repoRoot);
  printInstructions(repoRoot, skills);
  process.exit(0);
}

main();
