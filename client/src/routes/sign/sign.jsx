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
    padding: 2rem;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    align-items: center;
    background: linear-gradient(45deg, #0054d2, #00d2f7, #1c62d8);
   
    & .logo {
        max-width: 20rem;
        max-height: 20rem;
        border-radius: 50%;
        overflow: hidden;

        & img {
            width: 100%;
            height: 100%;
        }
    }

    & h1 {
        margin: 2rem;
        text-align: center;
        color: white;
        font-size: 1.5rem;
    }

    & .main-action {
        padding: 2rem 4rem;
        margin-bottom: 2rem;
        border-radius: 10px;
        box-shadow: var(--boxShadow);
        transition: .5s;

        &:hover{
            background: #ffffff2b;
        }

        & h3 {
            color: white;
            text-align: center;
            margin-bottom: 1rem;
        }
    }

    & .other-actions {
        display: flex;
        margin-bottom: 2rem;

        & span {
            margin: .5rem;
            padding: .5rem;
            border-radius: 5px;
            cursor: pointer;
            transition: var(--transitionApp);
            box-shadow: var(--boxShadow);

            &:hover {
                background: #002148;
            }
        }
    }

    & a {
        color: white;
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

const GInAction = ({ setAction }) => <span className="btn btn-primary" onClick={() => setAction("in")}>Войти</span>
const GUpAction = ({ setAction }) => <span className="btn btn-primary" onClick={() => setAction("up")}>Нет аккаунта? Зарегистроваться</span>
const GResetAction = ({ setAction }) => <span className="btn btn-primary" onClick={() => setAction("reset")}>Забыли пароль?</span>

export default function Sign() {
    const [action, setAction] = useState("in");

    return (
        <SSign>
            <div className="logo">
                <img src="/assets/app/logo.png" alt="logo" />
            </div>

            <h1 className="container-fluid">
                Al-Ber — это сервис, соединяющий человека, которому необходимо срочно отправить посылку, с тем, кому по пути
            </h1>

            <GSignAction action={action} setAction={setAction} />

            <Link to="/parsel" onClick={()=>USER.guest = true}>Войти как гость</Link>

            <Link target="_blank" download={true} to="/assets/rights/оферта.docx">Пользовательское соглашение и публичная оферта</Link>
        </SSign>
    )
}