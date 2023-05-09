from flask import Flask, request, jsonify
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

app = Flask(__name__)

CHROMEDRIVER = '/opt/chrome/chromedriver'

options = Options()
options.add_argument('--headless')  
options.add_argument('--no-sandbox')
options.add_argument('--disable-dev-shm-usage')

@app.route('/scrape', methods=['POST'])
def scrape():
    data = request.get_json()
    url = data.get('url', None)
    if not url:
        return jsonify({'error': 'No URL provided'}), 400

    chrome_service = ChromeService(executable_path=CHROMEDRIVER)
    driver = webdriver.Chrome(service=chrome_service, options=options)
    
    driver.get(url)
    
    # Scroll to the bottom of the page
    driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
    
    # Wait for the page to load
    WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.TAG_NAME, "body")))

    html = driver.page_source
    driver.quit()

    return jsonify({'html': html})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
