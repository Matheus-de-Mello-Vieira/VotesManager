# Paredão BBB

## Entendimento do problema

* Devo criar uma interface WEB para isso

* os usuários não precisam estar logados

* >  o usuário precisa receber um panorama percentual os votos por candidato até aquele momento

  * considerando o grande fluxo, se eu realizar a soma de todos os votos toda vez terei problemas de performance.

* >  o programa não quer receber votos oriundos de uma máquina, apenas votos de pessoas

  * colocar um CAPTCHA

* ele precisa ser elástico (ter autoscaling)

* `o total geral de votos, o total por participante e o total de votos por hora de cada paredão.`

  * diferente do cálculo para o usuário, esse é muito menos frequente e precisa de maior precisão, viabilizando calcular na hora

## Definições gerais

#### a parte do panorama percentual toda vez que alguém vota (candidatos de soluções):

1. **somar todos os votos toda vez que alguém votar**: como é uma rotina que seria invocada muito frequentemente, e eu precisaria de somar todos os votos que tem uma grande quantidade, realizar esses cálculos não seria algo viável.
   * Além disso, essas chamadas provavelmente obteriam resultados parecidos para obter resultados muito parecidos, de forma que seria um processamento desnecessário.
2. **armazenar o valor da soma total e atualizar toda vez que alguém votar**: como seria poucos registros (um para cada participante) sendo atualizados pro milhares de *threads* ao mesmo tempo, não seria algo viável devido a lentidão por *locks*.
3. **armazenar o valor da soma total e atualizar periodicamente**: acredito que funcione, pois assim a função de votar vai apenas adicionar um registro, e a função de somar vai apenar ler eles.
   * as *threads* responsáveis pela votação escrevem registros diferentes, então não teria necessidade de *locks* no banco de dados.
   * a *thread* responsável pela soma vai apenas ler os registros escritos pelas *threads* da votação, de forma a eliminar a necessidade de *locks*
     * isso inclusive me cheira o padrão **CQRS**
   * como o resultado que eu preciso mostrar para o usuário já está pronto, a votação pode ser registrada de maneira assíncrona.

Isso posto, eu vou adotar a solução 3, agora vou refinar ela:

1. o sistema não precisa somar tudo, podemos simplesmente manter no banco de dados o registro da hora da última soma, e apenas completar a soma já existe com os registros feitos apartir dessa data 
   * O problema disso seriam dois:
     * **performance**: toda vez que essa rotina for executada, ela fará um query no banco de dados para filtros todos os votos que aconteceram depois de uma data especifica. Se o SGBD não suporta ter timestamp enquanto índice primário, então esse solução é inviável.
       * considerando que os índices do bancos de dados são feitos através de uma árvore B ou uma variação dessa, normalmente apenas o índice primário é ordenado (chave primária ou qualquer estrutura parecida)
         * o fato da chave de busca (no caso a data) ter um índice ordenado é importante pois apenas nesses casos existe um fila encadeada que bastaria encontrar o primeiro elemento e retornar todos os que o segue. Os bancos de dados fazem isso automaticamente quando existe esse índice. 
     * **consistência**: seria necessário uma transação dentro do banco dados para manter a consistência, que é algo que preferia evitar por problemas de performance
       * aqui são os cenários caso ele não exista:
         * usar a data antes da soma: os votos que acontecerão entre esse meio termo serão somadas duas vezes: uma na soma atual, e outra na próxima, pois eles terão uma data superior que data salva no banco da última soma
         * usar a data depois da soma: os votos que acontecerão entre esse meio termo nunca serão contabilizados, pois foram adicionados depois que a soma foi feita, e na próxima soma também não serão excluíveis por ter uma data inferior a data salva no banco da última soma.
       * se for usar transação, ele precisa ser *serializable*, pois qualquer nível de isolamento inferior causaria inconsistência
         * considerando o nível *repeatable read* que é um nível abaixo do *serializable* por permitir obtenção de resultados parciais de operações feitas por outros agentes ativos no momento que a transação foi iniciada (*Phanton*): se dois votos foram feitos ao mesmo tempo no inicio da transação, um pode entrar na soma e o outro não (que se perderá para sempre pelos motivos descritos acima).
2. A mesma solução da 1, mas ao invés de data, usar uma chave primária serializado, de forma que a rotina somadora precisa apenas armazenar o id do último evento somado
   * Apesar de ter superar a solução 1, ele pressupõe que seria viável um banco de dados estabelecer mais de 1000 conexões por segundos, principalmente porque eu pretendo testar o sistema na minha máquina local, então não teria tantos containers.
   * Ao invés disso, poderíamos pensar em armazenar os dados agrupados por tempo e por candidato, ao invés de apenas armazenar os dados simples.
     * A forma como agruparia vou abordar no próximo tema, do sistema de votação ficar assíncrono.
     * mas esse agrupamento facilitaria também a parte do código para a produção do BBB, pois teria que somar muito menos dados

#### sistema assíncrono

1000 mensagens por segundo não é algo que se possa subestimar, e como eu quero testar o meu sistema na minha máquina, não teria como eu levantar 1000 threads aqui para rodar isso tudo. Ao invés disso, quando o usuário vota, o sistema vai apenas enfileirar isso.

O serviço na ponta consumidora vai baixar várias mensagens ao mesmo tempo (para economizar o número de conexões com o kafka) e vai agrupar isso na memória e salvar no banco de dados com uma conexão só.

Para evitar o uso de locks, e poder dar autoridade a memória local do sistema consumidor, eu gostaria que mensagens que seriam agrupadas no mesmo grupo não aparacerem em outras máquinas. Atualmente é possível implementar isso com o Kafka:

*  O Kafka particiona as mensagens de acordo com o atributo `partition key`, mensagens com esse atributo iguais sempre vão para a mesma partição.
* O Kafka garante que uma partição seja consumida apenas por um consumidor.
* Então se eu colocar um `partition key` relativo aos agrupamentos que quero implementar, eu terei tudo que preciso.

#### Qual arquitetura usar

Acredito que colocar todo o código em um serviço apenas não funcionaria pois tenho demandas com necessidade de escalonamento muito diferentes: um para lidar com os somatórios e outro para lidar com os registros das somas.

Os requisitos me pedem páginas HTML e a linguagem GO, de forma que a última que única opção que tenho é criar APIs HTTP usando o GO que responda as páginas em HTML.

Nesse sentido, eu penso em alguns serviços

* **votersInterface**:
  * interface do usuário comum
    * isso implica uma interface http que responda html
  * envia o registro de votação para a fila de votos
  * "view" aqui tem um sentido de "telespectador da TV"
* **votesRegister**:
  * consumer da fila de votos
  * ele registra no banco de dados

* **votesAggregator**:
  * invocado pelo cron
  * periodicamente atualizar a soma que vai para os usuários
  
* **productionInterface**:
  * interface da produção
  * executar as somas refinadas sob demanda da produção

Além disso, eu precisaria de:

* **um banco de dados**:
  * armazenar o histórico de votos
  * executar somas agrupadas
  * suportar querys mais complexas
* **um sistema de mensageria**:
  * desacoplar o recebimento do voto do se registro no banco de dados
  * evitar perda dos votos em momento de pico

#### Quais tecnologias usar

* Para lidar com isso, eu fiquei entre uma tecnologia tipo AWS Lambda, ou uma tecnologia de containers.
  * a tecnologia do tipo lambda tem melhor escalabilidade, mas é mais cara e mais difícil de testar
    * o problema disso é que seja quaisquer tecnologia que escolher, eu vou ficar preso a uma nuvem especifica (AWS, ou azure ou google)
* Eu vou usar o kubernetes para gerenciamento das containers
  * As principais núvens oferecem serviços para rodar o kubernetes, de forma que o sistema não ficaria limitado a um
  * com o kubernetes eu posso configurar tanto o banco de dados, quanto os containers das aplicações

* O banco de dados que vou usar vai ser o PostgreSQL
  * o PostgreSQL é o melhor banco de dados com transações open-source que conheço.
* O sistema de fila que vou usar vai ser o Kafka
  * os motivos falei acima, na parte do sistema assíncrono

#### Arquitetura interna do código

Apesar de ser um sistema relativamente simples, que por se só não justificaria uma arquitetura sofisticada como é o caso da arquitetura limpa, eu vou usar ela (parcialmente) como forma de demostração das minhas habilidades.

Por questão simplicidade, eu vou manter tudo no mesmo repositório, mas seria interessante imaginar isso como sendo um projeto grande, distribuído em nesses repositórios:

* `voting-commons`: códigos comuns entres os serviços, imagine isso enquanto uma biblioteca versionadada via CI/CD
  * vai conter os datamappers para comunicar com o Kafka (pois é usado em 2 componentes) e os datamappers para comunicar com o Postman (pois é usado por 3 componentes).

* `voters-frontend`: responsável pela interface web dos telespectadores
* `votes-register`: responsável por consumir a fila do Kafka e salvar os dados agrupados no banco de dados
* `votes-aggregator`: responsável por periodicamente pegar os dados registrados pelo **votes-register** e salvar a soma atualizada
* `prodution-frontend`: responsável pela interface dos telespectadores

## Coisas que implementaria se tivesse tempo

* No sistema existe uma relação no banco de dados agrupados por tempo, seria interessante começar agrupar os agrupamentos mais antigos, de forma a economizar espaço no disco. Poderia ser algo como: dados com mais de 2 semanas se agrupam por hora e com mais de 1 mês se agrupa por dia.
* Separa os repositórios e os módulos
  * incluir o `voting-commons` em algum repositório de artefatos, de forma que cada componente pode puxar uma versão diferente dele.
* imagens dos participantes
* criação de um pipeline CI/CD
  * SonarQube
* implementar o recaptcha
* Eu tentei usar o template para gerar os HTMLs, o problema é que isso levou a um código com baixa testabilidade, pois todos os dados ficam dentro do HTML.
  * ao invés disso, vou pegaria um arquivo HTML estático e vou preencher os elementos dele com javascript, chamando as rotas do meu sistema.
  * esse erro eu cometi apenas no `voter-frontend`, pois foi o primeiro que fiz. O `prodution-frontend` não terá o mesmo erro
* Pode personalizar o dashboard, agora ele é apenas por hora, mas talvez pensar ele em minutos, se

