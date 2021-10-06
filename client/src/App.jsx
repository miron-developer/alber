import { useEffect } from 'react';
import { useHistory } from 'react-router';

import { USER } from 'constants/constants';
import { IS_SIGN } from 'utils/content';
import Aside from 'components/aside/aside';
import Header from 'components/header/header';
import Main from 'components/routes/routes';

import './App.css';
import styled from 'styled-components';

const SApp = styled.div`
  	display: grid;
    grid-template-areas: ${props => props.isSign ? "'main'" : '"main aside" "header aside"'};
    grid-template-rows: ${props => props.isSign ? "1fr" : '1fr 7vh'};
    grid-template-columns: ${props => props.isSign ? '1fr' : '1fr 0fr'};
    min-height: 100vh;
    overflow: hidden;

	@media screen and (max-width: 600px) {
		& {
			grid-template-areas: '"main aside" "header aside"';
			grid-template-columns: 1fr;
			grid-template-rows: 1fr 0fr;
		}
	}
`;

export default function App() {
    const isSign = IS_SIGN();
    const history = useHistory();

    useEffect(() => {
        if (USER.status === "online") return isSign && history.push("/parsel")
        if (!USER.guest) return history.push('/sign');
    }, [history, isSign]);

    return (
        <SApp isSign={isSign}>
            {
                isSign
                    ? null
                    : <>
                        <Aside />
                        <Header />
                    </>
            }

            <Main />
        </SApp>
    )
};
