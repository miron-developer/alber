import { NavLink } from 'react-router-dom';

import styled from 'styled-components';

const SHeader = styled.header`
    grid-area: header;
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 1rem;
    width: 100vw;
    display: flex;
    justify-content: space-between;
    background: var(--blueColor);
    z-index: 5;
`;

const SNavLink = styled(NavLink)`
    margin: 0.5rem;
    padding: 0.5rem;
    width: calc(25% - 1rem);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid #231E2F;
    border-radius: 5px;
    color: var(--onHoverColor);
    text-decoration: none;
    transition: var(--transitionApp);

    &.active,
    &:hover {
        background: var(--onHoverColor);
        color: #192955;
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

export default function Header() {
    return (
        <SHeader>
            <GNavLink isExact={true} to="/parsel"       linkText="Отправит посылку" />
            <GNavLink isExact={true} to="/parsels"      linkText="Лента посылок" />
            <GNavLink isExact={true} to="/traveler"     linkText="Я попутчик" />
            <GNavLink isExact={true} to="/travelers"    linkText="Лента попутчиков" />
        </SHeader>
    )
}