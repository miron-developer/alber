import { Link } from "react-router-dom"

import styled from "styled-components"

const S404 = styled.section`
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;

    & h2 {
        font-size: 10rem;
        color: var(--redColor);
    }

    & h3 {
        font-size: 5rem;
        color: var(--darkRedColor);
    }

    & a {
        padding: 1rem;
        border-radius: 5px;
    }
`;

export default function Page404() {
    return (
        <S404>
            <h2>404</h2>
            <h3>Not found</h3>

            <Link to="/parsel">Go home</Link>
        </S404>
    )
}