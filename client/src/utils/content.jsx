import { useLocation } from "react-router";

import { Notify } from "components/app-notification/notification";

export const GET_CUR_PATHNAME = () => decodeURI(useLocation().pathname);
export const IS_SIGN = () => GET_CUR_PATHNAME().split('/').includes("sign");

// return two digit string
export const IsTwoDigit = data => parseInt(data) < 10 ? "0".concat(data) : data;

// generate date(15/12/2020 09:09) from milliseconds
export const DateFromMilliseconds = (milliseconds, isForInput = true) => {
    const datetime = new Date(parseInt(milliseconds));
    if (isNaN(datetime)) return;
    const str = [datetime.getFullYear(), IsTwoDigit(datetime.getMonth() + 1), IsTwoDigit(datetime.getDate())].join('-')
    return str;
}

// calculate time after thing created
export const CalculateRelativeDatetime = (datetime = Date.now().toString()) => {
    const now = Date.now();
    const given = parseInt(datetime);

    const delims = [[60, "мин."], [60, "ч."], [24, "д."], [7, 30, "нед."], [30, 365, "мес."], [365, "лет"]]
    const isPast = (now - given) > 0 ? true : false;
    let diff = Math.abs(now - given) / 1000;
    let diffType = "секунд";

    const returnDate = () => {
        diff = Math.floor(diff);

        if (!isPast) return ["через", IsTwoDigit(diff), diffType].join(' ');

        const isToday = Math.abs(now - given) < 86400000;
        return [isToday ? 'сегодня' : '', IsTwoDigit(diff), diffType, "назад"].join(' ');
    }

    // calculate difference
    for (let i = 0; i < delims.length; i++) {
        const cur_del = delims[i]
        if (cur_del.length === 2) {
            if (diff > cur_del[0]) {
                diff /= cur_del[0];
                diffType = cur_del[1];
                continue
            }
            break
        }
        if (cur_del[0] < diff && diff < cur_del[1]) {
            diff /= cur_del[0];
            diffType = cur_del[2];
            continue
        }
        break
    }
    return returnDate();
}

export const RandomKey = () => Math.round(Math.random() * Math.random() * 100000000);

export const ValidateParselTravelerSearch = (from, to, weight, price) => {
    // filter transit points, check exist, !=
    if (from === to || !to || !from) return Notify('fail', "Ошибка пунктов отправки и назначения")
    return {
        'fromID': from,
        'toID': to,
        'weight': weight || '',
        'price': price || '',
    }
}
