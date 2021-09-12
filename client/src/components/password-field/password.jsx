import { ShowAndHidePassword } from 'utils/effects';
import Input from '../form-input/input';

import styled from 'styled-components';

const SPasswordWrapper = styled.div`
    display: flex;
    align-items: baseline;
    cursor: pointer;

    & input {
        margin: 0 1rem;
        flex-grow: 1;
    }

    & i {
        padding: 5px;
        background: var(--offHoverBG);
        border-radius: 5px;
        box-shadow: 2px 2px 2px 0 #00000061;
    }
`;

export default function PasswordField({index, id, required, labelText, placeholder, pass, passToggle}) {
    return (
        <SPasswordWrapper>
            <Input index={index} id={id} type="password" name="password" base={pass.base} required={required}
                minLength="8" maxLength="30" placeholder={placeholder} labelText={labelText}
            />
            
            <i  className="fa fa-eye fa-eye-slash" 
                aria-hidden="true" 
                title="show/hide password"
                onClick={ e => {
                    ShowAndHidePassword(e, document.getElementById(id), passToggle);
                }} 
            ></i>
        </SPasswordWrapper>
    )
}