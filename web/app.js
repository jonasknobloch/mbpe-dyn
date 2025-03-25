const go = new Go();

document.addEventListener("DOMContentLoaded", function () {
    const tokenizeButton = document.getElementById("js-tokenize");
    const loadingMessage = document.getElementById("js-message");

    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(result => {
        go.run(result.instance);

        tokenizeButton.disabled = false;
        loadingMessage.style.display = "none";
    });

    const pinchZoom = document.getElementById("js-mermaid");

    pinchZoom.setTransform({
        scale: 1,
        x: 0,
        y: 220,
        allowChangeEvent: false,
    });

    tokenizeButton.addEventListener("click", () => tokenizeText());
});

let data = [];

async function tokenizeText() {
    const input = document.getElementById("js-input").value.trim();
    const model = document.getElementById("js-model").value;
    const vocab = parseInt(document.getElementById("js-vocab").value, 10) || -1;

    const message = document.getElementById("js-message");
    const mermaid = document.getElementById("js-mermaid");

    if (!input || typeof tokenizeWeb == "undefined") {
        return;
    }

    const result = JSON.parse(tokenizeWeb(input, model, vocab));

    if (!Array.isArray(result) || result.length < 1) {
        message.textContent = "No tokens found.";
        mermaid.innerHTML = "";

        return;
    }

    data = result;

    updateTabs(result);
    showResult(0);
}

function updateTabs() {
    const tabs = document.getElementById("js-tabs");

    tabs.innerHTML = "";

    data.forEach((chunk, index) => {
        const tab = document.createElement("button");

        tab.classList.add("mbpe-tab");
        tab.dataset.index = index.toString();
        tab.onclick = () => {
            tab.scrollIntoView({
                behavior: 'smooth',
                block: 'nearest',
                inline: 'center',
            });

            showResult(index);
        }

        const list = document.createElement("ol");

        list.classList.add("mbpe-chunk");

        chunk.Segmentations[chunk.Segmentations.length - 1].map(token => {
            const item = document.createElement("li");

            item.textContent = token.Token;
            item.classList.add("mbpe-token");
            item.title = token.ID;

            list.appendChild(item)
        });

        tab.appendChild(list);
        tabs.appendChild(tab);
    });
}

function showResult(index) {
    if (index < 0 || index >= data.length) {
        return;
    }

    mermaid.render("mermaidDiagram", data[index].Mermaid).then(({ svg }) => {
        document.getElementById("js-mermaid").innerHTML = svg;
    });

    document.querySelectorAll(".mbpe-tab").forEach(tab => tab.classList.remove("mbpe-tab--active"));
    document.querySelector(`.mbpe-tab[data-index="${index}"]`).classList.add("mbpe-tab--active");
}
