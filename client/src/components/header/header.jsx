import styled from 'styled-components';

const SHeader = styled.header`
    grid-area: header;
    padding: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--purpleColor);

    @media screen and (max-width: 600px) {
        position: fixed;
        left: 0;
        right: 0;
        height: 7vh;
        z-index: 5;
    }
`;

export default function Header() {
    return (
        <SHeader>
            
        </SHeader>
    )
}