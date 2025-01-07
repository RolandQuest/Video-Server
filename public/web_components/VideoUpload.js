


const template = document.createElement("template");
template.innerHTML = `
    <style>
        .box {
            
        }
        #drop_zone {
        
            background-color: slategray;
            
            border: 2px solid ivory;
            border-style: dashed;
            
            aspect-ratio: 3 / 2;
            width: 200px;
            
            text-align: center;
        }
        #drop_message {
            font-size: medium;
            color: snow;
            vertical-align: middle;
        }
    </style>
    
    <div class="box" id="drop_zone">
        <span id="drop_message"> Drop a file here! </span>
    </div>
`

export class VideoUpload extends HTMLElement {
    
    #shadow;
    #dropZone;
    #dropMessage;
    #debug;
    
    constructor() {
        super()
        
        this.#shadow = this.attachShadow({ mode: "open" })
        this.#shadow.append(template.content.cloneNode(true))
        
        this.#dropZone = this.#shadow.getElementById('drop_zone');
        this.#dropMessage = this.#shadow.getElementById('drop_message');
        
        this.#dropZone.addEventListener('drop', this);
        this.#dropZone.addEventListener('dragover', this);
        
        this.#debug = this.hasAttribute("debug");
    }
    
    handleEvent(ev) {
        ev.preventDefault();
        switch(ev.type) {
            case 'drop':
                this.#onDrop(ev);
                break;
            case 'dragover':
                this.#onDragOver(ev);
                break;
        }
    }
    
    #onDrop(ev) {
        const dropped_files = ev.dataTransfer.files;
        
        if(!dropped_files) {
            return;
        }
        if(dropped_files.length == 0) {
            return;
        }
        
        const file = dropped_files.item(0);
        
        // check file type is 'video/*'
        
        this.#dropMessage.innerText = `
            File Name: ${file.name}
            File Size: ${ (file.size / 1024 / 1024 / 1024).toLocaleString(undefined, { maximumFractionDigits: 2, minimumFractionDigits: 2 }) } GB
        `;
        
        let formData = new FormData();
        formData.append("video", file);
        fetch('/test/upload', {method: 'POST', body: formData })
        
    }
    
    #onDragOver(ev) {
        
    }
    
}

customElements.define("rq-video-upload", VideoUpload)
console.log('Web component imported: rq-video-upload')
