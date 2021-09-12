import { useState } from 'react';
import { useHistory } from 'react-router';

import { SubmitFormData, useInput } from 'utils/form';
import { Notify } from 'components/app-notification/notification';
import Input from 'components/form-input/input';
import SubmitBtn from 'components/submit-btn/submit';
import PhoneField from 'components/phone-field/field';

let afterStyles = []; // form handle all ::after notifications

export default function SignUp() {
    const nickname = useInput('');
    const phone = useInput('');
    const code = useInput('');
    const [step, setStep] = useState(1);
    const history = useHistory();

    const fields = [phone, nickname, code]; // fields for reset

    const onSuccessStep1 = () => {
        Notify('success', "Отправлено смс на номер " + phone.base.value + ". Возьмите оттуда код подтверждения")
        setStep(2);
    }
    const onSuccessStep2 = data => {
        Notify('success', `Вы успешно зарегистрированы. Ваш логин: "${phone.base.value}" и временный пароль:"${data?.password}"`, false)
        history.push("/")
    }
    const onFail = err => Notify('fail', 'Ошибка регистрации:' + err);

    return (
        step === 1
            ? <form action="/sign/sms/up" onSubmit={async (e) => {
                afterStyles = await SubmitFormData(e, afterStyles, fields, undefined, onSuccessStep1, onFail, false);
            }}>
                <h3>Регистрация(шаг 1)</h3>

                <PhoneField index="0" base={phone.base} />

                <SubmitBtn value="Отправить!" />
            </form>

            : <form action="/sign/up" onSubmit={async (e) => {
                afterStyles = await SubmitFormData(e, afterStyles, fields, undefined, onSuccessStep2, onFail);
            }}>
                <h3>Регистрация(шаг 2)</h3>

                <Input disabled={true} index="0" id="sign-up-phone" type="tel" name="phone" base={phone.base} labelText="Телефон:"/>
                <Input index="2" id="sign-up-code" type="text" name="code" base={code.base} labelText="8-значный код:"
                    minLength="8" maxLength="8" placeholder="Mfa7sd45"
                />
                <Input index="1" id="sign-up-nickname" type="text" name="nickname" base={nickname.base} labelText="Имя(никнейм):"
                    minLength="3" maxLength="20" placeholder="Miron"
                />
                <SubmitBtn value="Отправить!" />
            </form>

    )
}
