import { Redirect, useHistory } from 'react-router';
import { Switch, Route } from 'react-router-dom';

import { RandomKey } from 'utils/content';
import SignPage from 'routes/sign/sign';
import Page404 from 'routes/404/404';
import ParselPage from 'routes/parsel/parsel';
import ParselsPage from 'routes/parsels/parsels';
import TravelerPage from 'routes/traveler/traveler';
import TravelersPage from 'routes/travelers/travelers';
import FaqPage from 'routes/faq/faq';
import ContactsPage from 'routes/contacts/contacts';
import AdminPage from 'routes/admins/admins';
import Popup from 'components/popup/popup';
import AppNotifications from 'components/app-notification/notification';

import styled from 'styled-components';

const SMain = styled.main`
    grid-area: main;
    background: linear-gradient(45deg, #0054d2, #00d2f7, #1c62d8);

    & > * {
        margin-bottom: 5rem;
    }
`;

// app's routes
const ROUTES = [ {
    href: "/sign",
    isExact: false,
    component: SignPage,
}, {
    href: "/parsel",
    isExact: true,
    component: ParselPage,
}, {
    href: "/parsels",
    isExact: true,
    component: ParselsPage,
}, {
    href: "/traveler",
    isExact: true,
    component: TravelerPage,
}, {
    href: "/travelers",
    isExact: true,
    component: TravelersPage,
}, {
    href: "/faq",
    isExact: true,
    component: FaqPage,
}, {
    href: "/contacts",
    isExact: true,
    component: ContactsPage,
}, {
    href: "/admin",
    isExact: true,
    component: AdminPage,
}]

export default function DefineRoutes() {
    const history = useHistory();

    if (history.location.pathname === "/") return <Redirect to="/parsel" />
    return (
        <SMain>
            <Switch>
                {
                    ROUTES.map(
                        ({ href, component, isExact }) => <Route key={RandomKey()} exact={isExact} path={href} component={component} />
                    )
                }
                <Route component={Page404} />
            </Switch>

            <Popup />
            <AppNotifications />
        </SMain>
    )
}