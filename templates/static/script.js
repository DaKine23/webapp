var Root = React.createClass({
  getInitialState: function () {
    return {};
  },
  componentDidMount: function () {
    var component = this;
    $.get("/pong", function (data) {
      component.setState(data)
    });
  },
  render: function () {
    return (
      <h1>{this.state.message}</h1>
    )
  }
});

ReactDOM.render(
  <Root />,
  document.getElementById('root')
);
