import libtorrent as lt
import tempfile
from flask import Flask, request, send_file

app = Flask(__name__)

@app.route('/download', methods=['POST'])
def download_torrent():
    url = request.form.get('url')
    ses = lt.session()
    info = lt.torrent_info(url)
    h = ses.add_torrent({'ti': info, 'save_path': '/tmp'})
    ses.start_dht()

    while not h.is_seed():
        s = h.status()
        print('Downloading:', s.progress, '%')

    return send_file('/tmp/' + info.name() + '.torrent', as_attachment=True)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
