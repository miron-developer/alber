import { useHistory } from 'react-router';

import { USER } from 'constants/constants';
import { SubmitFormData, useInput, useTogglePassword } from 'utils/form';
import { UserOnline } from 'utils/user';
import { Notify } from 'components/app-notification/notification';
import Input from 'components/form-input/input';
import PasswordField from 'components/password-field/password';
import SubmitBtn from 'components/submit-btn/submit';

let afterStyles = []; // form handle all ::after notifications

export default function SignIn() {
    const login = useInput();
    const pass = useInput();
    const passToggle = useTogglePassword();
    const fields = [login, pass];
    const history = useHistory()

    // custom validation
    const customValidation = () =>
        (!(/[a-z]+/g.test(pass.base.value) && /[A-Z]+/g.test(pass.base.value) && /[0-9]+/g.test(pass.base.value))) ? [1] : [];

    const onSuccess = async (data) => {
        const isOnline = await UserOnline(data.id);
        if (isOnline) {
            Notify('success', "Вход произведен");
            USER.guest = false;
            history.push('/parsel');
        } else Notify('fail', "Ошибка входа")
    }
    const onFail = err => Notify('fail', "Ошибка входа:" + err);

    return (
        <form action="/sign/in" onSubmit={async (e) => {
            afterStyles = await SubmitFormData(e, afterStyles, fields, customValidation, onSuccess, onFail);
        }}>
            <h3>Вход</h3>

            <Input index="0" id="login" name="phone" type="tel" base={login.base} labelText="Логин:"
                minLength="11" maxLength="15" placeholder="87777777777"
            />
            <PasswordField index="1" id="password" name="password" labelText="Пароль:"
                placeholder="User1234" pass={pass} passToggle={passToggle}
            />

            <SubmitBtn value="Отправить!" />
        </form>
    )
}