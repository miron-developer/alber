import { IS_SIGN } from '@functions/content';
import Aside from '@common/aside/aside';
import Header from '@common/header/header';
import Footer from '@common/footer/footer';
import Main from '@ommon/routes/routes';

import './App.css';
import styled, { css } from 'styled-components';

const generalStyle = css`
    :root {
        --onHoverColor: #FFFFFF;
        --offHoverColor: #D3D3D3;
        --purpleColor: #6B5B95;
        --violetColor: #2F0B8D;
        --redColor: #FF0000;
        --darkRedColor: rgb(102, 12, 12);
        --successBG: #02ff0075;
        --failBG: #ff00007a;
        --infoBG: #f9ff007d
        --asideBG: linear-gradient(180deg, #6B5B95 0%, #CDABB2 100%);
        --navsBG: linear-gradient(180deg, #FFFFFF 0%, rgba(0, 0, 0, 0) 100%);
        --mainBG: linear-gradient(180deg, rgba(135, 40, 60, 0.69) 0%, rgba(189, 53, 82, 0) 100%);
        --onHoverBG: #E7C7CE;
        --offHoverBG: #CDABB2;
        --boxShadow: 5px 5px 10px rgb(0 0 0 / 25%);
        --transitionApp: .5s;
    }

    * {
        box-sizing: border-box;
    }

    body {
        margin: 0;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
    }

    code {
        font-family: source-code-pro, Menlo, Monaco, Consolas, 'Courier New', monospace;
    }

    @media screen and (max-width: 900px) {
        * {
            font-size: 12px;
        }
    }

    @media screen and (max-width: 700px) {
        * {
            font-size: 11px;
        }
    }

    @media screen and (max-width: 600px) {
        * {
            font-size: 10px;
        }
    }
`

const SApp = styled.div`
    ${generalStyle}

  	display: grid;
    grid-template-areas: ${props => props.isSign ? '"header" "main" "footer"' : '"aside header" "aside main" "footer footer"'};
    grid-template-rows: 7vh 1fr 7vh;
    grid-template-columns: ${props => props.isSign ? '1fr' : '2fr 8fr'};
    min-height: 100vh;

	@media screen and (max-width: 600px) {
		& {
			grid-template-areas: "aside" "header" "main" "footer";
			grid-template-columns: 1fr;
			grid-template-rows: 0 7vh 1fr 7vh;
		}
	}
`;

export default function App() {
    const isSign = IS_SIGN();

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
            <Footer />
        </SApp>
    )
};
