import styled from 'styled-components';

const SFormField = styled.div`
    margin-bottom: .5rem;
`;

const SFormInputWrapper = styled.div`
    display: flex;
    align-items: center;
    margin-bottom: .5rem;
`;

const SFormInputLabel = styled.label`
    white-space: nowrap;
    color: var(--offHoverColor);

    &.required::after {
        content: '*';
        color: var(--redColor);
    }
`;

const SFormInput = styled.input`
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

    &::placeholder{
        color: var(--darkGreyColor);
    }

    &:disabled {
        background: var(--greyColor);
    }
`;

const SFormInputNotification = styled.div`
    color: var(--darkRedColor);
`;

export const Label = ({ required, id, labelText }) =>
    <SFormInputLabel className={required ? 'required' : ''} htmlFor={id} > {labelText} </SFormInputLabel>

export default function Input({ index, id, type = "text", name, labelText, base, minLength, maxLength, min, max, required = true, hidden = false, placeholder = "" }) {
    return hidden
        ? <input type={type} value={base.value} name={name} hidden />
        : (
            <SFormField className={'form-field-' + index}>
                <SFormInputWrapper>
                    <Label required={required} id={id} labelText={labelText} />
                    <SFormInput
                        className="form-input"
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
                        {...base}
                    />
                </SFormInputWrapper>
                <SFormInputNotification className="form-input-notification" />
            </SFormField>
        )
}