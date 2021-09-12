import { useState } from 'react';
import { Link } from 'react-router-dom';

import { USER } from 'constants/constants';
import { RandomKey } from 'utils/content';

import SignUp from './up/up';
import SignIn from './in/in';
import ResetPassword from './reset/reset';
import styled from 'styled-components';

const SSign = styled.section`
    height: 100%;
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    align-items: center;

    & .logo {
        max-width: 40vw;
        max-height: 40vh;

        & img {
            width: 100%;
            height: 100%;
        }
    }

    & .main-action {
        min-width: 60vw;
        padding: 2rem;
        border-radius: 10px;
        background: #fdfdfd;
        box-shadow: var(--boxShadow);
    }

    & .other-actions {
        display: flex;
        flex-direction: column;

        & span {
            margin: .5rem;
            padding: .5rem;
            border-radius: 5px;
            cursor: pointer;
            transition: var(--transitionApp);

            &:hover {
                background: var(--onHoverBG);
                color:  var(--onHoverColor);
            }
        }
    }
`;

const GSignAction = ({ action, setAction }) => {
    let mainAction;
    let otherActions = []
    if (action === "up") {
        mainAction = <SignUp />;
        otherActions = [<GInAction key={RandomKey()} setAction={setAction} />, <GResetAction key={RandomKey()} setAction={setAction} />];
    } else if (action === "reset") {
        mainAction = <ResetPassword />;
        otherActions = [<GUpAction key={RandomKey()} setAction={setAction} />, <GInAction key={RandomKey()} setAction={setAction} />];
    } else {
        mainAction = <SignIn />;
        otherActions = [<GUpAction key={RandomKey()} setAction={setAction} />, <GResetAction key={RandomKey()} setAction={setAction} />];
    }
    return (
        <>
            <div className="main-action">{mainAction}</div>
            <div className="other-actions">
                {otherActions}
            </div>
        </>
    )
}

const GInAction = ({ setAction }) => <span onClick={() => setAction("in")}>Войти</span>
const GUpAction = ({ setAction }) => <span onClick={() => setAction("up")}>Нет аккаунта? Зарегистроваться</span>
const GResetAction = ({ setAction }) => <span onClick={() => setAction("reset")}>Забыли пароль?</span>

export default function Sign() {
    const [action, setAction] = useState("in");

    return (
        <SSign>
            <div className="logo">
                <img src="/assets/app/logo512.png" alt="logo" />
            </div>

            <GSignAction action={action} setAction={setAction} />

            <Link to="/" onClick={()=>USER.guest = true}>Войти как гость</Link>

            <Link download to="/assets/name.txt">Пользовательское соглашение</Link>
        </SSign>
    )
}