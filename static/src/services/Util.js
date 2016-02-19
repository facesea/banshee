/*@ngInject*/
module.exports = function () {
  return {
    buildRepr: function (rule) {
      rule.trendUp = rule.trendUp || false;
      rule.trendDown = rule.trendDown || false;
      rule.thresholdMax = rule.thresholdMax || 0;
      rule.thresholdMin = rule.thresholdMin || 0;

      var parts = [];
      if (rule.trendUp && rule.thresholdMax === 0) {
        parts.push('trend ↑');
      }
      if (rule.trendUp && rule.thresholdMax !== 0) {
        parts.push('(trend ↑ && value >= ' + parseFloat(rule.thresholdMax.toFixed(3)) + ')');
      }
      if (!rule.trendUp && rule.thresholdMax !== 0) {
        parts.push('value >= ' + parseFloat(rule.thresholdMax.toFixed(3)));
      }
      if (rule.trendDown && rule.thresholdMin === 0) {
        parts.push('trend ↓');
      }
      if (rule.trendDown && rule.thresholdMin !== 0) {
        parts.push('(trend ↓ && value <= ' + parseFloat(rule.thresholdMin.toFixed(3)) + ')');
      }
      if (!rule.trendDown && rule.thresholdMin !== 0) {
        parts.push('value <= ' + parseFloat(rule.thresholdMin.toFixed(3)));
      }
      return parts.join(' || ');
    }
  };
};
