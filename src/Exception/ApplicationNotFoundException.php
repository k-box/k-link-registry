<?php

namespace App\Exception;

class ApplicationNotFoundException extends KRegistryException
{
    /**
     * @var string[]
     */
    private $errors;

    public function __construct(array $errors = [])
    {
        $this->errors = $errors;
    }

    /**
     * @return string[]
     */
    public function getErrors(): array
    {
        return $this->errors;
    }
}
