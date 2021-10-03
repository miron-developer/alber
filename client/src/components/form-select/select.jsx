import { RandomKey } from "utils/content";

import styled from "styled-components";

const SLalel = styled.label`
    display: flex;
    align-items: center;
    white-space: nowrap;

    & span {
        color: white;
    }

    & select {
        margin-left: 1rem;
    }
`;
export default function Select({ name, text, required = true, options, onChange }) {
    return (
        <SLalel>
            <span>{text}</span>
            <select className="form-select" name={name} value={options.selected} required={required} onChange={onChange}>
                {options?.data?.map((opt) => <option key={RandomKey()} value={opt[options.value]}>{options.makeText(opt)}</option>)}
            </select>
        </SLalel>
    )
}