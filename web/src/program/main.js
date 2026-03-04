const inputPanel  = document.getElementById('inputPanel');
const resultPanel = document.getElementById('resultPanel');
const urlInput    = document.getElementById('urlInput');
const aliasInput  = document.getElementById('aliasInput');
const expirySelect = document.getElementById('expirySelect');
const urlError    = document.getElementById('urlError');
const shortenBtn  = document.getElementById('shortenBtn');
const shortUrlDisplay = document.getElementById('shortUrlDisplay');
const origUrlDisplay  = document.getElementById('origUrlDisplay');
const expiryDisplay   = document.getElementById('expiryDisplay');
const copyBtn   = document.getElementById('copyBtn');
const copyLabel = document.getElementById('copyLabel');
const backBtn = document.getElementById('backBtn');
const backBtn2 = document.getElementById(`backBtn2`)
const BACKEND_URL = import.meta.env.VITE_BACKEND_URL

function isValidUrl(str) {
  try { const u = new URL(str); return u.protocol === 'http:' || u.protocol === 'https:'; }
  catch { return false; }
}

function showPanel(show, hide) {
  hide.classList.remove('visible-panel');
  hide.classList.add('hidden-panel');
  show.classList.remove('hidden-panel');
  show.classList.add('visible-panel');
}

async function handleShorten() {
  const url = urlInput.value.trim();

  if (!isValidUrl(url)) {
    urlError.classList.remove('hidden');
    urlInput.classList.add('!border-red-300');
    urlInput.focus();
    return;
  }
  
  urlError.classList.add('hidden');
  urlInput.classList.remove('!border-red-300');

  // Loading
  shortenBtn.innerHTML = `<span class="dot-loader"><span></span><span></span><span></span></span>`;
  shortenBtn.disabled = true;

  const expiryMap = {'1d': 1, '3d': 3, '7d' : 7, '30d': 30 };

  // ── Replace this await with your actual fetch() to your backend ──
  const res = await fetch(`${BACKEND_URL}/shorten`,{
    method: "POST",
    headers: {"Content-type":"application/json"},
    body: JSON.stringify({
      url: url,
      alias: aliasInput.value.trim() || null,
      expires: expiryMap[expirySelect.value],
    }),
  });

  if (!res.ok) {
    shortenBtn.disabled = false;
    console.error(`Failed fetching: ${res.status}`)
    return
  }
  
  console.info(res.status)
  const data = await res.json();
  console.log("> ", data)

  const shortUrl = `localhost:8080/${data.ShortCode}`; // change this to domain name after deployment

  shortUrlDisplay.textContent = shortUrl;
  origUrlDisplay.textContent  = url;
  expiryDisplay.textContent = `${expiryMap[expirySelect.value]} days`;

  shortenBtn.innerHTML = `
    <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.8" class="w-4 h-4" stroke-linecap="round">
      <path d="M10.5 3.5l3 3-3 3M9.5 16.5l-3-3 3-3"/>
      <path d="M13.5 6.5H8a4 4 0 000 8h1"/>
    </svg>
    Shorten URL`;
  shortenBtn.disabled = false;

  showPanel(resultPanel, inputPanel);
}

function handleBack() {
  showPanel(inputPanel, resultPanel);
  urlInput.value     = '';
  aliasInput.value   = '';
  expirySelect.value = '1d';
  copyLabel.textContent = 'Copy';
  copyBtn.classList.remove('copied');
}

function handleCopy() {
  navigator.clipboard.writeText(shortUrlDisplay.textContent).then(() => {
    copyBtn.classList.add('copied');
    copyLabel.textContent = '✓ Copied';
    setTimeout(() => {
      copyBtn.classList.remove('copied');
      copyLabel.textContent = 'Copy';
    }, 2000);
  });
}

urlInput.addEventListener('keydown', e => { if (e.key === 'Enter') handleShorten(); });
shortenBtn.addEventListener('click', handleShorten);
copyBtn.addEventListener('click', handleCopy);
backBtn.addEventListener('click', handleBack);
backBtn2.addEventListener('click', handleBack);