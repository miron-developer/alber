import { useState } from 'react';
import { NavLink, useHistory } from 'react-router-dom';

import { USER } from 'constants/constants';
import { SignOut } from 'utils/user';

import styled from 'styled-components';

const SAside = styled.aside`
    grid-area: aside;
    padding: 1rem;
    max-width: 30vw;

    @media screen and (max-width: 600px) {
        position: fixed;
        right: -100vw;
        height: 100vh;
        width: 80vw;
        max-width: 80vw;
        z-index: 10;
        opacity: .9;
        transition: calc(var(--transitionApp)*2);
        
        &.open {
            transform: translate(-100vw);
        }
    }
`

const SAsideTop = styled.div`
    margin: 1rem;
`

const SLogo = styled.div`
    margin: auto;
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
    margin: .5rem auto;
    padding: .5rem;
    width: max-content;
    max-width: 100%;
    text-transform: uppercase;
    color: var(--purpleColor);
    font-weight: bold;
    text-align: center;
    word-break: break-all;
    background: var(--onHoverColor);
    border-radius: 5px;
    transition: .5s;
`

const SLogout = styled(SNickname)`
    color: var(--redColor);
    cursor: pointer;
    transition: var(--transitionApp);

    &:hover {
        background: var(--redColor);
        color: var(--onHoverColor);
    }
`

const SNavs = styled.nav`
    padding: 1rem;
    display: flex;
    flex-direction: column;
    background: var(--navsBG);
`

const SNavLink = styled(NavLink)`
    margin: 0.5rem 0;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--purpleColor);
    border: 1px solid #231E2F;
    border-radius: 5px;
    color: var(--onHoverColor);
    text-shadow: 1px 1px 5px black;
    text-decoration: none;
    text-transform: uppercase;
    transition: var(--transitionApp);

    &.active,
    &:hover {
        color: var(--purpleColor);
        text-shadow: none;
        background: var(--onHoverColor);
    }
`

// Generate navlink
const GNavLink = ({isExact, to, linkText}) => {
    return (
        <SNavLink exact={isExact} activeClassName="active" to={to}>
            <span className="nav-link-text">{linkText}</span>
        </SNavLink>
    )
}

export default function Aside() {
    const [isOpened, setOpened] = useState(false);
    const history = useHistory();

    return (
        <SAside className="aside" onClick={()=> setOpened(!isOpened)}>
            <SAsideTop>
                <SLogo as={NavLink} to="/" >
                    <img src="/img/logo192.png" alt="wnet logo" />
                </SLogo>

                <SAsideTop>
                    <SNickname>{USER.nickname}</SNickname>
                    <SLogout onClick={() => SignOut(history)}></SLogout>
                </SAsideTop>
            </SAsideTop>

            <SNavs>
                <GNavLink
                    isExact={true}
                    to="/" 
                    linkText=""
                />
            </SNavs>

        </SAside>
    )
}