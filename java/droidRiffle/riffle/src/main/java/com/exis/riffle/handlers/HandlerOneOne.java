package com.exis.riffle.handlers;

/**
 * Created by damouse on 1/24/2016.
 */

public interface HandlerOneOne<A, R> extends AnyHandler {
    R run(A a);
}