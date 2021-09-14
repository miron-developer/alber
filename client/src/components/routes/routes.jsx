import { Switch, Route } from 'react-router-dom';

import { RandomKey } from 'utils/content';
import SignPage from 'routes/sign/sign';
import Page404 from 'routes/404/404';
import ParselPage from 'routes/parsel/parsel';
import ParselsPage from 'routes/parsels/parsels';
import TravelerPage from 'routes/traveler/traveler';
import TravelersPage from 'routes/travelers/travelers';
import Popup from 'components/popup/popup';
import AppNotifications from 'components/app-notification/notification';

import styled from 'styled-components';

const SMain = styled.main`
    grid-area: main;
    background: var(--greyColor);
`;

// app's routes
const ROUTES = [{
    href: "/sign",
    isExact: false,
    component: SignPage,
}, {
    href: "/parsel",
    isExact: true,
    component: ParselPage,
},{
    href: "/",
    isExact: true,
    component: ParselPage,
},{
    href: "/parsels",
    isExact: true,
    component: ParselsPage,
},{
    href: "/traveler",
    isExact: true,
    component: TravelerPage,
}, {
    href: "/travelers",
    isExact: true,
    component: TravelersPage,
}]

export default function DefineRoutes() {
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

            {/* popups */}
            <Popup />
            <AppNotifications />
        </SMain>
    )
}