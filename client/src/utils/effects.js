import { Notify } from "components/app-notification/notification";
import { PopupOpen } from "components/popup/popup";
import ManageParsel from "components/parsel/manage/manage";
import ManageTraveler from "components/traveler/manage/manage";
import ToTopType from "components/to-toptype/toptype";
import ToUp from "components/to-up/up";


import { GetDataByCrieteries, POSTRequestWithParams } from "./api";
// just debounce
export function Debounce(fn, ms) {
    let timeOut;
    return (...args) => {
        clearTimeout(timeOut);
        timeOut = setTimeout(() => { fn(...args) }, ms)
    }
}

export const DbnceCities = Debounce(async(e) => {
    const res = await GetDataByCrieteries("search", { 'type': 'cities', 'q': e.target.value });
    if (res.err) return Notify("fail", "Не удалось загрузить города");

    const options = res.map(({ id, name }) => {
        const op = document.createElement("option");
        op.value = name;
        op.textContent = id;
        return op;
    })

    const dt = document.getElementById(e.target.list.id);
    if (!dt) return;
    dt.innerHTML = "";
    dt.append(...options);
}, 500)

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
 * @param childrenClass get parent by childs classnames
 * @param isStopLoad stop load or no
 * @param isScrollingToTop load on scroll to top or bottom
 * @param loadCallback what do after react edge
 */
export const ScrollHandler = Debounce(async(e, isStopLoad, isScrollingToTop = false, loadCallback = async() => {}, childrenClass) => {
    if (isStopLoad) return;

    const parent = e.target;

    console.log('par', parent);
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

export const EditItem = async(type, data, cb) =>
    PopupOpen(type === "parsel" ? ManageParsel : ManageTraveler, { 'cb': cb, 'data': data, 'type': 'edit' })

export const RemoveItem = async(id, type, cb) => {
    const res = await POSTRequestWithParams("/r/" + (type === "parsel" ? "parsel" : 'travel'), { 'id': id })
    if (res.err && res.err !== "ok") return Notify('fail', 'Не удалено');
    cb()
}

export const TopItem = async(id, type, cb) => PopupOpen(ToUp, { 'cb': cb, "type": type, 'id': id })

export const PaintItem = async(id, type, cb) => PopupOpen(ToTopType, { 'cb': cb, "type": type, 'id': id })

const removeEmptyFields = (obj = {}) => {
    for (let [k, v] of Object.entries(obj))
        if (v === "" || !v) delete obj[k];
    return obj
}

export const CompareParams = (newParams, currentParams) => {
    const res = {};
    newParams = removeEmptyFields(newParams);
    for (let [k, v] of Object.entries(newParams)) {
        if (newParams[k] !== currentParams[k]) {
            res[k] = v;
        }
    }
    return res;
}

export const GetValueFromListByIDAndInputValue = (listID, inputValue) => {
    const dt = Array.from(document.getElementById(listID).childNodes)
    if (dt.length === 0) return;
    const op = dt.find(option => option.value.includes(inputValue));
    if (op) return op.textContent;
}

export const OnChangeTransitPoint = async(point, e) => {
    point.setCertainValue(e.target.value);
    DbnceCities(e);
}