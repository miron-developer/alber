import { Switch, Route } from 'react-router-dom';

import { RandomKey } from 'utils/content';
import HomePage from 'routes/home/home';
import SignPage from 'routes/sign/sign';
import Page404 from 'routes/404/404';
import Popup from 'components/popup/popup';
import AppNotifications from 'components/app-notification/notification';

import styled from 'styled-components';

const SMain = styled.main`
    grid-area: main;
    background: var(--greyColor);
`;

// app's routes
const ROUTES = [{
    href: "/",
    isExact: true,
    component: HomePage,
}, {
    href: "/sign",
    isExact: false,
    component: SignPage,
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