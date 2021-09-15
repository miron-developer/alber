import styled from "styled-components"

const SSubmitBtn = styled.input`
    width: 60%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 1rem auto;
    padding: 1rem;
    font-size: 1rem;
    color: var(--offHoverColor);
    background: var(--offHoverBG);
    border-radius: 5px;
    border: none;
    box-shadow: 2px 2px 2px 0 #00000061;
    transition: var(--transitionApp);
    cursor: pointer;

    &:hover {
        color: var(--onHoverColor);
        background: var(--onHoverBG);
    }
`;

export default function SubmitBtn({value, onClick}) {
    return <SSubmitBtn type="submit" value={value} onClick={onClick} />
}