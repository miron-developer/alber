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
    width: 80vw;
    height: 100vh;
    background: #4e4a4a66;
    transition: var(--transitionApp);
    z-index: 10;
    opacity: .9;

    &.open {
        transform: translate(-80vw);
    }

    & .open-btn {
        position: absolute;
        right: 100%;
        top: 80%;
        padding: 1rem;
        border-radius: 5px;
        font-size: 1.5rem;
        color: white;
        background:var(--blueColor);
        z-index: 15;
        cursor: pointer;
    }
`

const SUser = styled.div`
    margin: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;

    & > * {
        display: flex;
        align-items: center;
    }
`

const SLogo = styled.div`
    margin: auto;
    width: 10vw;
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
    word-break: break-all;
    background: var(--onHoverColor);
    border-radius: 5px;
    transition: var(--transitionApp);
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
    const history = useHistory();

    return (
        <SAside className={isOpened ? "open" : ""}>
            {/* aside open/close btn */}
            <div className="open-btn" onClick={() => setOpened(!isOpened)}>
                <i className="fa fa-bars" aria-hidden="true"></i>
            </div>

            {
                isOpened &&
                <>
                    <SLogo as={NavLink} to="/" >
                        <img src="/assets/app/logo192.png" alt="logo" />
                    </SLogo>

                    <SUser>
                        {
                            USER.status === "online"
                                ? <div>
                                    <SNickname>
                                        <i className="fa fa-user" aria-hidden="true"></i>
                                        {USER.nickname} ({USER.phoneNumber})
                                    </SNickname>
                                    <SEdit onClick={()=>PopupOpen(EditProfile, {})}>
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

                    {
                        USER.status === "online" && <History />
                    }
                </>
            }
        </SAside>
    )
}