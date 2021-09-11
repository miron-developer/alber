// just debounce
export function Debounce(fn, ms) {
    let timeOut;
    return (...args) => {
        clearTimeout(timeOut);
        timeOut = setTimeout(() => { fn(...args) }, ms)
    }
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