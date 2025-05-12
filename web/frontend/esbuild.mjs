import * as esbuild from 'esbuild';
import fg from 'fast-glob';
import path from 'path';
import fs from 'fs';

fs.rmSync('../assets/dist', { recursive: true });
const timestamp = Date.now();

const jsFiles = await fg('../assets/js/**/*.js');
const cssFiles = await fg('../assets/css/**/*.css');

const allFiles = [...jsFiles, ...cssFiles];
const entries = allFiles.map((file) => {
  const parsed = path.parse(file);
  return {
    in: file,
    out: `${parsed.name}.${timestamp}${parsed.ext}`,
  };
});

await esbuild.build({
  entryPoints: entries.map((e) => e.in),
  entryNames: `[name].${timestamp}`,
  format: 'cjs',
  bundle: true,
  minify: true,
  outdir: '../assets/dist',
  loader: {
    '.css': 'css',
  },
});

const manifest = Object.fromEntries(
  entries.map((e) => [path.basename(e.in), path.basename(e.out)]),
);

fs.writeFileSync(
  path.join('../assets/dist', 'manifest.json'),
  JSON.stringify(manifest, null, 2),
);
