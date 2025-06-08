export async function loadModule(url: string) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to load module from ${url}`);
  }
  const scriptText = await response.text();
  // const blob = new Blob([scriptText], { type: 'module' });
  const blob = new Blob([scriptText], { type: 'application/javascript' });
  const blobUrl = URL.createObjectURL(blob);
  return import(/* webpackIgnore: true */ blobUrl);
}
