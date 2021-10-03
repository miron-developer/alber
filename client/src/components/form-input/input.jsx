import styled from 'styled-components';

const SFormField = styled.div`
    margin-bottom: .5rem;
    display: flex;

    & label {
        /* margin: .5rem; */
        /* color: white; */
        /* white-space: nowrap; */
    }

    & input {
        box-shadow: 4px 4px 3px 0 #00000029;
    }
`;

const SFormInputLabel = styled.label`
    white-space: nowrap;
    color: var(--offHoverColor);

    &.required::after {
        content: '*';
        color: var(--redColor);
    }
`;

export const Label = ({ required, id, labelText }) =>
    <SFormInputLabel className={required ? 'required' : ''} htmlFor={id} > {labelText} </SFormInputLabel>

export default function Input({ index, id, type = "text", name, labelText, base, minLength, maxLength, list, min, max, required = true, hidden = false, placeholder = "" }) {
    return hidden
        ? <input type={type} value={base.value} name={name} hidden />
        : (
            <SFormField className={'form-floating form-field-' + index}>
                <input
                    className="form-control"
                    id={id}
                    type={type}
                    name={name}
                    required={required}
                    min={min}
                    max={max}
                    minLength={minLength}
                    maxLength={maxLength}
                    placeholder={placeholder}
                    hidden={hidden}
                    list={list}
                    {...base}
                />
                <label htmlFor={id} placeholder={labelText}>{labelText}</label>
            </SFormField>
        )
}