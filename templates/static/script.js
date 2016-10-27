var Root = React.createClass({
  getInitialState: function () {
    return {};
  },
  componentDidMount: function () {
    var component = this;
    doPoll(component);
  },
  render: function () {
    return (
      <h1>{this.state.message}</h1>
    )
  }
});
function doPoll(component){
    $.get("/pong", function (data) {
      component.setState(data)
      i = i + 1
      if (i<2)
      setTimeout(function () {doPoll(component)},200);
      i = i - 1
    });
};

var i = 0;

ReactDOM.render(
  <Root />,
  document.getElementById('root')
);
