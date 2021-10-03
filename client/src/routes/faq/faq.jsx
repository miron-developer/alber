import styled from "styled-components"

const SFaq = styled.article`
    padding: 2rem;

    & > * {
        margin: 2rem 0;
    }

    & h1.faq__title {
        color: white;
    }

    & p.faq__intro {
        font-size: 1.2rem;
    }

    & .accordion-body{
        background: #192955;
        color: white;
    }
`;

export default function FaqPage() {
    return (
        <SFaq className="faq">
            <h1 className="faq__title">Вопросы и ответы</h1>
            <p className="faq__intro">
                Здесь собраны ответы на большинство вопросов, которые интересуют наших клиентов. Пожалуйста, ознакомьтесь с ними. В случае, если вы не найдете ответа на свой вопрос, пожалуйста, свяжитесь с нами.
            </p>

            <div className="accordion" id="faq-accordion">
                <div className="accordion-item">
                    <h2 className="accordion-header" id="safety">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
                            А это безопасно?
                        </button>
                    </h2>
                    <div id="collapseOne" className="accordion-collapse collapse show">
                        <div className="accordion-body">
                            Если соблюдать все правила, да. Мы настоятельно рекомендуем проверять содержимое передаваемых (или доставляемых) посылок. Для нас важна безопасность наших пользователей, поэтому мы разработали список обязательных к соблюдению <a className="faq-list__link" href="/assets/rights/name.txt" target="_blank">правил</a>.
                        </div>
                    </div>
                </div>
                <div className="accordion-item">
                    <h2 className="accordion-header" id="payment">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseTwo" aria-expanded="true" aria-controls="collapseTwo">
                            Сколько я заработаю?
                        </button>
                    </h2>
                    <div id="collapseTwo" className="accordion-collapse collapse">
                        <div className="accordion-body">
                            Al-Ber позволяет заработать попутчикам. Если вы собираетесь в поездку, вы можете взять с собой посылку или документы. Доставить их в оговоренное место (на вокзал, в аэропорт или в городе) и получить за это вознаграждение.
                        </div>
                    </div>
                </div>
                <div className="accordion-item">
                    <h2 className="accordion-header" id="responsibility">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree" aria-expanded="true" aria-controls="collapseThree">
                            Какая ответственность за доставку?
                        </button>
                    </h2>
                    <div id="collapseThree" className="accordion-collapse collapse">
                        <div className="accordion-body">
                            Перед использованием нашего сервиса рекомендуем ознакомиться с <a className="faq-list__link" href="/assets/rights/name.txt" target="_blank">правилами</a>. Сервис не несет ответственность за сохранность документов или посылок. Al-Ber помогает найти попутчика и заработать на своей поездке.
                        </div>
                    </div>
                </div>
                <div className="accordion-item">
                    <h2 className="accordion-header" id="max-weight">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFour" aria-expanded="true" aria-controls="collapseFour">
                            Посылку какого максимального веса я могу отправить?
                        </button>
                    </h2>
                    <div id="collapseFour" className="accordion-collapse collapse">
                        <div className="accordion-body">
                            Как правило, с помощью Al-Ber отправляют документы или посылки до 5 кг. В исключительных случаях (когда попутчик едет на поезде или имеет багаж 20 кг) можно отправить посылку с большим весом.
                        </div>
                    </div>
                </div>
                <div className="accordion-item">
                    <h2 className="accordion-header" id="price">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFive" aria-expanded="true" aria-controls="collapseFive">
                            Сколько стоит отправить посылку или документ?
                        </button>
                    </h2>
                    <div id="collapseFive" className="accordion-collapse collapse">
                        <div className="accordion-body">
                            Преимущества краудшиппингового сервиса в том, что стоимость отправления устанавливаете вы. Мы не диктуем цены, лишь можем подсказать рыночную стоимость по вашему направлению.
                        </div>
                    </div>
                </div>
            </div>
        </SFaq>
    )
}