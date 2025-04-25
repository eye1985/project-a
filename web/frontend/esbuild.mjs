import * as esbuild from 'esbuild';

await esbuild.build({
  entryPoints: [
    { in: '../assets/dist/index.js', out: `index.${new Date().getTime()}` },
    { in: '../assets/dist/shortcut.js', out: `shortcut.${new Date().getTime()}` }
  ],
  bundle: true,
  minify: true,
  outdir: '../assets/dist'
});