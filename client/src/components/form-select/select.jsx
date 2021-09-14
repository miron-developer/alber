import { RandomKey } from "utils/content";

import styled from "styled-components";

const SLalel = styled.label`
    display: flex;
    align-items: center;
    white-space: nowrap;

    & select {
        margin-left: 1rem;
        padding: .5rem;
        width: 100%;
        font-size: 1rem;
        color: var(--offHoverColor);
        background: none;
        border: none;
        border-radius: 5px;
        outline: none;
        border-bottom: 1px solid var(--onHoverColor);
        box-shadow: 4px 4px 3px 0 #00000029;
    }
`;
export default function Select({ name, text, required, options, onChange }) {
    return (
        <SLalel>
            <span>{text}</span>
            <select name={name} required={required} onChange={onChange}>
                {options?.data?.map((opt) => <option key={RandomKey()} value={opt[options.value]}>{options.makeText(opt)}</option>)}
            </select>
        </SLalel>
    )
}