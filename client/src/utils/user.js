import { USER } from 'constants/constants';
import { GetDataByCrieteries, POSTRequestWithParams } from './api';
import { Notify } from 'components/app-notification/notification';

export const IsLogged = async() => {
    const res = await POSTRequestWithParams('/sign/status');
    if (res.err !== 'ok') {
        Notify('fail', res.err);
        return false;
    }
    return res.data.id;
}

// change USER const
const changeUserData = async(id) => {
    if (id !== undefined) {
        const res = (await GetDataByCrieteries('user', { id: id }));
        if (res.err && res.err !== "ok") {
            return Notify('fail', "Ошибка входа в аккаунт. Не верные данные")
        }
        for (let [k, v] of Object.entries(res[0])) USER[k] = v;
        USER.status = 'online';
        return true;
    } else {
        for (let k in USER) USER[k] = '';
        USER.status = Date.now().toString();
        USER.nickname = "nickname";
        return true;
    }
}

// Switch to online
export const UserOnline = async(id) => await changeUserData(id);;

// Switch to offline
export const UserOffline = async() => await changeUserData();

// send to server signal about sign out
export const SignOut = async(history) => {
    Notify('info', "Производится выход...");
    const res = await POSTRequestWithParams("/sign/out");
    if (res.err !== "ok") return Notify('fail', "Ошибка: выход не произведен");;
    const isSignOuted = await UserOffline();
    if (!isSignOuted) return;
    history.push('/sign/in');
    Notify('success', "Выход произведен");
}