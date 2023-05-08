const express = require('express');
const puppeteer = require('puppeteer');

const app = express();

app.use(express.urlencoded({ extended: true }));
app.use(express.json());

app.post('/scroll', async (req, res) => {
  const { url } = req.body;

  const browser = await puppeteer.launch({
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  const page = await browser.newPage();
  await page.goto(url, { waitUntil: 'networkidle2' });

  await autoScroll(page);

  const html = await page.content();
  res.send(html);

  await browser.close();
});

async function autoScroll(page) {
  await page.evaluate(async () => {
    await new Promise((resolve, reject) => {
      let totalHeight = 0;
      const distance = 100;
      const timer = setInterval(() => {
        const scrollHeight = document.body.scrollHeight;
        window.scrollBy(0, distance);
        totalHeight += distance;

        if (totalHeight >= scrollHeight) {
          clearInterval(timer);
          resolve();
        }
      }, 100);
    });
  });
}

app.listen(process.env.PORT || 3000, () => {
  console.log('Server is running');
});
