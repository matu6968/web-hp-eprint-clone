const puppeteer = require('puppeteer');
const path = require('path');
const { URL } = require('url');

// Check if a file path or URL is provided as an argument
const inputPath = process.argv[2];
if (!inputPath) {
  console.error('Please provide the path to the HTML file or a URL as an argument.');
  process.exit(1);
}

// Determine if the input is a URL
let isUrl;
try {
  new URL(inputPath);
  isUrl = true;
} catch (error) {
  isUrl = false;
}

// Set the output file name (optional argument)
const outputFile = process.argv[3] || 'output.png';

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  // Set viewport size and device scale factor for better readability
  await page.setViewport({
    width: 1240,           // Set width of the viewport
    height: 1754,           // Set height of the viewport
    deviceScaleFactor: 2   // Increase the device scale factor for a zoom effect
  });

  // Load the page, either from a URL or a local file
  if (isUrl) {
    await page.goto(inputPath, { waitUntil: 'networkidle0' });
  } else {
    const fullPath = path.resolve(inputPath);
    await page.goto(`file://${fullPath}`, { waitUntil: 'networkidle0' });
  }

  // Capture the screenshot
  await page.screenshot({
    path: outputFile,
    fullPage: true
  });

  await browser.close();
  console.log(`Screenshot taken and saved as ${outputFile}`);
})();

