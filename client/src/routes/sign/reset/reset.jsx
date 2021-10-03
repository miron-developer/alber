import { useState } from 'react';
import { useHistory } from 'react-router';

import { SubmitFormData, useInput, useTogglePassword } from 'utils/form';
import { Notify } from 'components/app-notification/notification';
import Input from 'components/form-input/input';
import SubmitBtn from 'components/submit-btn/submit';
import PasswordField from 'components/password-field/password';
import PhoneField from 'components/phone-field/field';

let afterStyles = []; // form handle all ::after notifications

export default function Restore() {
    const phone = useInput();
    const code = useInput('');
    const pass = useInput('');
    const passToggle = useTogglePassword()
    const [step, setStep] = useState(1);
    const history = useHistory();
    const fields = [phone, pass, code]; // fields for reset

    const onSuccessStep1 = () => {
        Notify('success', "Отправлено смс на номер " + phone.base.value + ". Возьмите оттуда код подтверждения");
        setStep(2);
    }
    const onSuccessStep2 = () => {
        Notify('success', "Пароль успешно изменен.");
        history.push("/sign")
    }
    const onFail = err => Notify('fail', 'Ошибка восстановления:' + err);

    return step === 1
        ? <form action="/sign/sms/ch" onSubmit={async (e) => {
            afterStyles = await SubmitFormData(e, afterStyles, fields, undefined, onSuccessStep1, onFail, false);
        }}>
            <h3>Восстановление пароля (шаг 1)</h3>

            <PhoneField index="0" base={phone.base} />

            <SubmitBtn value="Отправить!" />
        </form>

        : <form action="/sign/re" onSubmit={async (e) => {
            afterStyles = await SubmitFormData(e, afterStyles, fields, undefined, onSuccessStep2, onFail);
        }}>
            <h3>Восстановление пароля (шаг 2)</h3>

            <PasswordField index="1" id="password" name="password" labelText="Новый пароль:"
                placeholder="User1234" pass={pass} passToggle={passToggle}
            />

            <Input index="2" id="code" type="text" name="code" base={code.base} labelText="6-значный код:"
                minLength="6" maxLength="6" placeholder="123456"
            />
            <SubmitBtn value="Отправить!" />
        </form>
}