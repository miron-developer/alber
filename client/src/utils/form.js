import { useState } from 'react';

import { Fetching } from './api';

// useInput custom react hook
export const useInput = (initValue = '') => {
    const [value, setValue] = useState(initValue);
    const handleChange = e => setValue(e.target.value);
    const resetField = () => setValue('');
    const setCertainValue = v => setValue(v);

    return {
        base: {
            value,
            onChange: handleChange,
        },
        resetField,
        setCertainValue
    }
}

// for show/hide password
export const useTogglePassword = (initState = 'password') => {
    const [state, changeType] = useState(initState);
    const toggleType = () => changeType(state === 'password' ? 'text' : 'password');

    return {
        state,
        toggleType
    }
}

// check xss attack
export const isHaveXSS = content => /<+[\w\s/]+>+/gm.test(content);

// base validation check html validation & xss
export const BaseFormValidation = data => {
    const resIndexs = [];
    let i = 0;
    for (let v of data.values()) {
        if (isHaveXSS(v)) resIndexs.push(i);
        i++;
    }
    return resIndexs;
}

// create one ::after notification
export const NotifyInvalidField = (selector, field) => {
    const afterStyle = document.createElement('style');
    afterStyle.innerHTML += `${selector}::after{content: 'Ошибка'; font-size: 1rem; color: red; margin: 10px;}`;
    document.head.appendChild(afterStyle);
    field.resetField();
    return afterStyle;
}

// reset form by calling component method
export const ResetForm = fields => {
    if (!fields) return;
    fields.forEach(field => field.resetField());
}

// clear all custom ::after styles
export const ClearAllNotifications = (afterStyles = []) => afterStyles.forEach(AS => AS.remove());

const setNotificationIfExist = (BIndexs, CIndexs, afterStyles, fields) => {
    if (!fields) return;
    if (BIndexs.length !== 0)
        BIndexs.forEach(index => afterStyles.push(NotifyInvalidField(`.form-field-${1+index} .form-input-notification`, fields[index])));
    if (CIndexs && CIndexs.length !== 0)
        CIndexs.forEach(index => afterStyles.push(NotifyInvalidField(`.form-field-${1+index} .form-input-notification`, fields[index])));
}

// submit form: check all validations, return ::after styles, notify & reset fields
export const SubmitFormData = async(e, afterStyles, fields, customValidation = () => {}, onSuccess = () => {}, onFail = () => {}, isReset = true) => {
    e.preventDefault();
    ClearAllNotifications(afterStyles);
    afterStyles = [];

    const form = e.target;
    const action = form.getAttribute('action');
    if (!action) return;
    const data = new FormData(form);

    // validation
    const BIndexs = BaseFormValidation(data),
        CIndexs = customValidation();

    // notification
    setNotificationIfExist(BIndexs, CIndexs, afterStyles, fields);

    // fetching & notification
    if ((BIndexs.length === 0 && !CIndexs) ||
        (BIndexs.length === 0 && CIndexs && CIndexs.length === 0)) {
        const res = await Fetching(action, data);

        if (res.err === "ok") {
            if (isReset) ResetForm(fields);
            return onSuccess(res.data);
        }
        return onFail(res.err);
    }
    onFail("Ошибка валидации");
    return afterStyles;
}