import { Switch, Route } from 'react-router-dom';

import NF from 'nf404/error';
import Home from 'home/home';
import Signs from 'signs/sign';
import Profile from 'profile/profile';
import Popup from 'common/popup/popup';

import styled from 'styled-components';

const SMain = styled.main`
    grid-area: main;
    background: var(--mainBG);
`;

// app's routes
const ROUTES = [{
        href: "/",
        isExact: true,
        component: Home,
    },
    {
        href: "/sign/",
        isExact: false,
        component: Signs,
    },
    {
        href: "/user/:id",
        isExact: true,
        component: Profile,
    },
]

export default function DefineRoutes() {
    return (
        <SMain>
            <Switch>
                {
                    ROUTES.map(
                        ({href, component, isExact }, index) => <Route key={index} exact={isExact} path={href} component={component} />
                    )
                }
                <Route component={NF} />
            </Switch>

            <Popup />
        </SMain>
    )
}