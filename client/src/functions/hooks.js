import { useState } from "react";

import { Notify } from "@common/app-notification/notification";

import { GetDataByCrieteries } from "./api";

// for lazy load
export const useFromTo = (initState = [], step = 10) => {
    const [fromToState, setFromToState] = useState({
        'start': 0,
        'isStopLoad': false,
        'datalist': initState,
    });

    const setDataList = state => setFromToState(Object.assign({}, fromToState, { 'datalist': state }));

    const getPart = async(getWhat = "", params = {}, failText = "", isAppToEnd = true, isNeedClear = false) => {
        if (getWhat === "" || failText === "") return Notify('fail', failText);

        const res = await GetDataByCrieteries(getWhat, {
            ...params,
            'from': fromToState.start,
            'step': step
        });

        if (res.err && res.err !== 'ok') {
            fromToState.isStopLoad = true;
            setFromToState(Object.assign({}, fromToState));
            return Notify('fail', failText + " : " + res.err);
        }

        if (isNeedClear) fromToState.datalist = res;
        else if (isAppToEnd) fromToState.datalist = [...fromToState.datalist, ...res];
        else fromToState.datalist = [...res, ...fromToState.datalist];

        if (res.length < step) fromToState.isStopLoad = true;
        else fromToState.start += step;

        setFromToState(Object.assign({}, fromToState));
        return true;
    }

    const zeroState = () => setFromToState(Object.assign({}, fromToState, {
        'start': 0,
        'isStopLoad': false,
        'datalist': initState,
    }));

    return {
        'datalist': fromToState.datalist,
        'isStopLoad': fromToState.isStopLoad,
        setDataList,
        getPart,
        zeroState,
    }
}