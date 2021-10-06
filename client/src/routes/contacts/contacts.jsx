import styled from "styled-components"

const SContacts = styled.article`
    padding: 2rem;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    color: white;

    & > * {
        margin: 1rem 0;
    }

    p.contacts__intro {
        font-size: 1.2rem;
    }

    section.contacts-section {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
    }

    .contacts-block {
        margin: 1rem;
        max-width: 50%;
        min-width: 25%;

        & a {
            text-decoration: none;
            display: flex;
            align-items: center;
            justify-content: center;
            flex-direction: column;
            color: #192955;
        }
    }

    .contacts-block i {
        font-size: 3rem;
        color: #1b68d9;
    }
`;

export default function ContactsPage() {
    return (
        <SContacts className="contacts">
            <h1 className="contacts__title">Контактная информация</h1>
            <p className="contacts__intro">
                Мы рады, что вы используете сервис Al-Ber.<br />
                Остались вопросы или пожелания по&nbsp; улучшению нашего сервиса — напишите нам.
            </p>

            <section className="contacts-section">
                <div className="contacts-block">
                    <a className="contacts-block__link" rel="noreferrer" href="https://wa.me/+77787833831" alt="" target="_blank">
                        <i className="fa fa-whatsapp" aria-hidden="true"></i>

                        <div className="contacts-block__right">
                            <span className="contacts-block__text">WhatsApp</span>
                        </div>
                    </a>
                </div>

                <div className="contacts-block">
                    <a className="contacts-block__link" href="mailto:Aibek.burkitbay@gmail.com">
                        <i className="fa fa-envelope" aria-hidden="true"></i>

                        <div className="contacts-block__right">
                            <span className="contacts-block__text">Электронная почта</span>
                        </div>
                    </a>
                </div>

                <div className="contacts-block">
                    <a className="contacts-block__link" href="tel:+77787833831">
                        <i className="fa fa-phone" aria-hidden="true"></i>

                        <div className="contacts-block__right">
                            <span className="contacts-block__text">+7 778 783-38-31</span>
                        </div>
                    </a>
                </div>

                <div className="contacts-block">
                    <a className="contacts-block__link" rel="noreferrer" href="https://www.instagram.com/Al_ber.kz" target="_blank">
                        <i className="fa fa-instagram" aria-hidden="true"></i>

                        <div className="contacts-block__right">
                            <span className="contacts-block__text">Instagram</span>
                        </div>
                    </a>
                </div>
            </section>
        </SContacts>
    )
}