import { migratePatchFunc, migratePatchAndWarnFunc } from '../main.js';

// Support jQuery slim which excludes the deferred module in jQuery 4.0+
if (jQuery.Deferred) {
  var oldDeferred = jQuery.Deferred,
    tuples = [
      // Action, add listener, callbacks, .then handlers, final state
      [
        'resolve',
        'done',
        jQuery.Callbacks('once memory'),
        jQuery.Callbacks('once memory'),
        'resolved',
      ],
      [
        'reject',
        'fail',
        jQuery.Callbacks('once memory'),
        jQuery.Callbacks('once memory'),
        'rejected',
      ],
      [
        'notify',
        'progress',
        jQuery.Callbacks('memory'),
        jQuery.Callbacks('memory'),
      ],
    ];

  migratePatchFunc(
    jQuery,
    'Deferred',
    function (func) {
      var deferred = oldDeferred(),
        promise = deferred.promise();

      function newDeferredPipe(/* fnDone, fnFail, fnProgress */) {
        var fns = arguments;

        return jQuery
          .Deferred(function (newDefer) {
            jQuery.each(tuples, function (i, tuple) {
              var fn = typeof fns[i] === 'function' && fns[i];

              // Deferred.done(function() { bind to newDefer or newDefer.resolve })
              // deferred.fail(function() { bind to newDefer or newDefer.reject })
              // deferred.progress(function() { bind to newDefer or newDefer.notify })
              deferred[tuple[1]](function () {
                var returned = fn && fn.apply(this, arguments);
                if (returned && typeof returned.promise === 'function') {
                  returned
                    .promise()
                    .done(newDefer.resolve)
                    .fail(newDefer.reject)
                    .progress(newDefer.notify);
                } else {
                  newDefer[tuple[0] + 'With'](
                    this === promise ? newDefer.promise() : this,
                    fn ? [returned] : arguments
                  );
                }
              });
            });
            fns = null;
          })
          .promise();
      }

      migratePatchAndWarnFunc(
        deferred,
        'pipe',
        newDeferredPipe,
        'deferred-pipe',
        'deferred.pipe() is deprecated'
      );
      migratePatchAndWarnFunc(
        promise,
        'pipe',
        newDeferredPipe,
        'deferred-pipe',
        'deferred.pipe() is deprecated'
      );

      if (func) {
        func.call(deferred, deferred);
      }

      return deferred;
    },
    'deferred-pipe'
  );

  // Preserve handler of uncaught exceptions in promise chains
  jQuery.Deferred.exceptionHook = oldDeferred.exceptionHook;
}
