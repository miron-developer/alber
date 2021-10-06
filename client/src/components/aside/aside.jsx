import { useState } from 'react';
import { NavLink, useHistory } from 'react-router-dom';

import { USER } from 'constants/constants';
import { SignOut } from 'utils/user';
import { PopupOpen } from 'components/popup/popup';

import History from './history/history';
import EditProfile from './edit-profile/profile';
import styled from 'styled-components';

const SAside = styled.aside`
    grid-area: aside;
    position: fixed;
    left: 100vw;
    padding: 1rem;
    width: 85vw;
    height: 100vh;
    background: #2b2b2be0;
    transition: var(--transitionApp);
    z-index: 10;
    opacity: .9;

    &.open {
        transform: translate(-85vw);
    }

    & .open-btn {
        position: absolute;
        right: 100%;
        top: 75%;
        padding: 1rem;
        border-radius: 5px;
        font-size: 1.5rem;
        color: white;
        background:var(--blueColor);
        z-index: 15;
        cursor: pointer;
    }

    & .menu {
        height: 65%;

        .menu-items > span {
            margin: .5rem;
            padding: 1rem;
            color: var(--onHoverColor);
            border-radius: 10px;
            transition: var(--transitionApp);
            cursor: pointer;

            &.active,
            &:hover {
                background: var(--blueColor);
            }
        }

        .links {
            display: flex;
            flex-direction: column;
            margin: 2rem 0;
        }

        .links a {
            padding: 1rem;
            margin-bottom: 1rem;
            font-weight: bold;
            text-decoration: none;
            background: white;
            border-radius: 10px;
            transition: .5s;

            &:hover,
            &.active {
                background: #193162;
                color: white;
            }
        }
    }
`

const SUser = styled.div`
    margin: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 20%;

    & > * {
        display: flex;
        align-items: center;
    }
`

const SLogo = styled.div`
    margin: auto;
    width: 10rem;
    height: 15%;
    display: block;
    overflow: hidden;
    transition: var(--transitionApp);

    &:hover {
        filter: brightness(0.5);
    }

    & img {
        height: 100%;
        width: 100%;
    }
`

const SNickname = styled.div`
    margin: .5rem;
    padding: .5rem;
    width: max-content;
    max-width: 100%;
    text-transform: uppercase;
    font-weight: bold;
    text-align: center;
    word-break: break-word;
    background: var(--onHoverColor);
    border-radius: 5px;
    transition: var(--transitionApp);

    & > i {
        margin: 5px;
    }

    &.nickname {
        width: 100%;
    }
`

const SLogout = styled(SNickname)`
    color: var(--redColor);
    cursor: pointer;

    &:hover {
        background: var(--redColor);
        color: var(--onHoverColor);
    }
`

const SEdit = styled(SNickname)`
    cursor: pointer;

    &:hover {
        background: var(--redColor);
        color: var(--onHoverColor);
    }
`;

export default function Aside() {
    const [isOpened, setOpened] = useState(false);
    const [tab, setTab] = useState(0);
    const history = useHistory();

    return (
        <SAside className={isOpened ? "open" : ""}>
            <div className="open-btn" onClick={() => setOpened(!isOpened)}>
                <i className="fa fa-bars" aria-hidden="true"></i>
            </div>

            {
                isOpened &&
                <>
                    <SLogo as={NavLink} to="/" >
                        <img src="/assets/app/logo.png" alt="logo" />
                    </SLogo>

                    <SUser>
                        {
                            USER.status === "online"
                                ? <div>
                                    <SNickname className="nickname">
                                        <i className="fa fa-user" aria-hidden="true"></i>
                                        <span>{USER.nickname} ({USER.phoneNumber})</span>
                                    </SNickname>
                                    <SEdit onClick={() => PopupOpen(EditProfile, {})}>
                                        <i className="fa fa-pencil" aria-hidden="true"></i>
                                    </SEdit>
                                </div>
                                : <SNickname>Здесь будет ваше имя</SNickname>
                        }

                        {
                            USER.status === "online"
                                ? <SLogout onClick={() => SignOut(history)}>
                                    <i className="fa fa-sign-out" aria-hidden="true"></i>
                                    Выход
                                </SLogout>
                                : <SLogout onClick={() => history.push("/sign")}>
                                    <i className="fa fa-sign-in" aria-hidden="true"></i>
                                    Войти
                                </SLogout>
                        }
                    </SUser>

                    <div className="menu">
                        <div className="menu-items">
                            <span className={tab === 0 ? 'active' : ''} onClick={() => setTab(0)}>Меню</span>

                            {USER.status === "online" && <span className={tab === 1 ? 'active' : ''} onClick={() => setTab(1)}>Ваша история</span>}

                        </div>
                        {
                            USER.status === "online" && tab === 1
                                ? <History />
                                : <div className="links">
                                    <NavLink to="/faq">Вопросы и ответы</NavLink>
                                    <NavLink to="/contacts">Контакты</NavLink>
                                    {USER.isAdmin && <NavLink to="/admin">Админ</NavLink>}
                                </div>
                        }
                    </div>
                </>
            }
        </SAside>
    )
}