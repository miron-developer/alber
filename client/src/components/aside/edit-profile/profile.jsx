import { useState } from "react";

import { USER } from "constants/constants";
import { SubmitFormData, useInput, useTogglePassword } from "utils/form";
import { UserOnline } from "utils/user";
import { Notify } from "components/app-notification/notification";
import PhoneField from "components/phone-field/field";
import SubmitBtn from "components/submit-btn/submit";
import Input from "components/form-input/input";
import PasswordField from "components/password-field/password";

import styled from "styled-components";

let afterStyles = [];

const SForms = styled.form `
    padding: 2rem;
`;

export default function EditProfile() {
    const nickname = useInput('');
    const phone = useInput('');
    const pass = useInput('');
    const passToggle = useTogglePassword()
    const code = useInput('');
    const [step, setStep] = useState(1);

    const fields = [phone, nickname, pass, code]; // fields for reset

    const onSuccessStep1 = () => {
        Notify('success', "Отправлено смс на номер " + phone.base.value + ". Возьмите оттуда код подтверждения")
        setStep(2);
    }
    const onSuccessStep2 = () => Notify('success', `Вы успешно изменили ваши данные.`) || UserOnline(USER.id)
    const onFail = err => Notify('fail', 'Ошибка регистрации:' + err);

    return (
        step === 1
            ? <SForms action="/e/user/confirm" onSubmit={async (e) => {
                afterStyles = await SubmitFormData(e, afterStyles, fields, undefined, onSuccessStep1, onFail, false);
            }}>
                <h3>Смена данных(шаг 1). Введите только то, что хотите</h3>

                <PhoneField index="0" base={phone.base} required={false} />

                <Input index="1" id="nickname" type="text" name="nickname" base={nickname.base} labelText="Имя(никнейм):"
                    minLength="3" maxLength="20" placeholder="Miron" required={false}
                />
                <PasswordField index="2" id="password" name="password" labelText="Пароль:"
                    placeholder="User1234" pass={pass} passToggle={passToggle} required={false}
                />

                <SubmitBtn value="Отправить!" />
            </SForms>

            : <SForms action="/e/user" onSubmit={async (e) => {
                afterStyles = await SubmitFormData(e, afterStyles, fields, undefined, onSuccessStep2, onFail);
            }}>
                <h3>Смена данных(шаг 2)</h3>

                <input hidden type="tel" name="phone" {...phone.base} />
                <input hidden type="text" name="nickname" {...nickname.base} />
                <input hidden type="password" {...pass.base} />
        
                <Input index="3" id="code" type="text" name="code" base={code.base} labelText="8-значный код:"
                    minLength="8" maxLength="8" placeholder="Mfa7sd45"
                />

                <SubmitBtn value="Отправить!" />
            </SForms>
    )
}