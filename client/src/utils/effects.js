import { Notify } from "components/app-notification/notification";
import { PopupOpen } from "components/popup/popup";
import { POSTRequestWithParams } from "./api";

// just debounce
export function Debounce(fn, ms) {
    let timeOut;
    return (...args) => {
        clearTimeout(timeOut);
        timeOut = setTimeout(() => { fn(...args) }, ms)
    }
}

// show & hide password by changing input type
export const ShowAndHidePassword = (e, passElem, passwordToggle) => {
    const elem = e.target;
    passwordToggle.toggleType();
    elem.classList.toggle('fa-eye-slash');
    if (passwordToggle.state === "password") passElem.setAttribute('type', 'text');
    else passElem.setAttribute('type', 'password');
}

/**
 * for lazy load and keeping focus with scrolling
 * @param e event
 * @param isStopLoad stop load or no
 * @param className for getting priorEdgeChild
 * @param isScrollingToTop load on scroll to top or bottom
 * @param loadCallback what do after react edge
 */
export const ScrollHandler = Debounce(async(e, isStopLoad, isScrollingToTop = false, loadCallback = async() => {}) => {
    if (isStopLoad) return;

    const parent = e.target;
    const pRec = parent.getBoundingClientRect();
    if (
        (isScrollingToTop && parent.scrollTop === 0) ||
        (!isScrollingToTop && parent.scrollTop >= Math.round((parent.scrollHeight - pRec.height) * .75))
    ) {
        const priorEdgeChildNum = isScrollingToTop ? 0 : parent.childElementCount - 1;

        if (await loadCallback()) {
            setTimeout(() => {
                // smooth scroll
                const el = parent.childNodes[priorEdgeChildNum];
                if (el) el.scrollIntoView({ behavior: "smooth" });
            }, 100);
        }
    }
}, 100);

export const EditItem = async(type, data, cb) => PopupOpen('edit', { 'cb': cb, "type": type, ...data })

export const RemoveItem = async(id, type, cb) => {
    const res = await POSTRequestWithParams("/r", { 'type': type, 'id': id })
    if (res.err && res.err !== "ok") return Notify('fail', 'Не удалено');
    cb()
}

export const TopItem = async(id, type, cb) => PopupOpen('toptype', { 'cb': cb, "type": type, 'id': id })

export const PaintItem = async(id, type, cb) => PopupOpen('up', { 'cb': cb, "type": type, 'id': id })