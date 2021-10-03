import { useCallback, useState } from "react";

import { Notify } from "components/app-notification/notification";

import { GetDataByCrieteries } from "./api";

// for lazy load
export const useFromTo = (initState = [], step = 10) => {
    const [fromToState, setFromToState] = useState({
        'start': 0,
        'isStopLoad': false,
        'isLoaded': false,
        'datalist': initState,
    });

    const setDataList = state => setFromToState(Object.assign({}, fromToState, { 'datalist': state }));

    const getPart = useCallback(async(getWhat = "", params = {}, failText = "", isAppToEnd = true, isNeedClear = false) => {
        if (getWhat === "" || failText === "") return Notify('fail', failText);

        const res = await GetDataByCrieteries(getWhat, {
            ...params,
            'from': isNeedClear ? 0 : fromToState.start,
            'step': step
        });

        if (res.err && res.err !== 'ok') {
            fromToState.isStopLoad = true;
            fromToState.isLoaded = true;
            setFromToState(Object.assign({}, fromToState));
            return Notify('fail', failText + " : " + res.err);
        }

        if (!fromToState.isLoaded) fromToState.isLoaded = true;
        if (isNeedClear) {
            fromToState.start = 0;
            fromToState.datalist = res;
            fromToState.isStopLoad = false;
        } else if (isAppToEnd) fromToState.datalist = [...fromToState.datalist, ...res];
        else fromToState.datalist = [...res, ...fromToState.datalist];

        if (res.length < step) fromToState.isStopLoad = true;
        else fromToState.start += step;

        setFromToState(Object.assign({}, fromToState));
        return fromToState.datalist;
    }, [fromToState, step])

    const zeroState = async() => setFromToState({
        'start': 0,
        'isStopLoad': false,
        'isLoaded': false,
        'datalist': initState,
    });

    return {
        'datalist': fromToState.datalist,
        'isStopLoad': fromToState.isStopLoad,
        'isLoaded': fromToState.isLoaded,
        setDataList,
        getPart,
        zeroState,
    }
}