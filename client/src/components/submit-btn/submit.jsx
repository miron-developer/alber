import styled from "styled-components"

const SSubmitBtn = styled.input`
    width: 60%;
    display: block;
    margin: 1rem auto;
    padding: 1rem;
    border-radius: 5px;
    border: none;
    box-shadow: 2px 2px 2px 0 #00000061;
    transition: var(--transitionApp);
    cursor: pointer;

    &:hover {
        background: #002148;
    }
`;

export default function SubmitBtn({value, onClick}) {
    return <SSubmitBtn className="btn btn-primary" type="submit" value={value} onClick={onClick} />
}