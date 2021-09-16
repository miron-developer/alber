import { POSTRequestWithParams } from "utils/api";

// preload file
export const PreloadFile = (accept, cb = (file, src, type) => {}) => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = accept;
    input.click();

    input.addEventListener('change', async(e) => {
        e.stopImmediatePropagation();
        const file = input.files[0];

        if (file) cb(file, URL.createObjectURL(file), file.type);
    });
}

// upload file to server
export const UploadFile = async(type, file, whomFile, whomID) => {
    if (!type || !file || !whomFile) return { err: 'deficite data' };
    const params = {
        'type': type,
        'file': file,
        'whomFile': whomFile,
    }
    if (whomID) params['whomID'] = whomID;

    const res = await POSTRequestWithParams('/s/image', params);
    if (res.err !== "ok") return null;
    return res.data;
}